package station

import (
	"fmt"

	"github.com/murtaza-u/antistrofi/ping"
	pb "github.com/murtaza-u/antistrofi/proto/gen/go"
	"github.com/murtaza-u/antistrofi/token"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var secret []byte

// service is the default service. It is always registered with the grpc
// server in addition to your custom services.
type service struct {
	pb.UnimplementedControlServer
}

func (s *service) Callback(_ *pb.Void, stream pb.Control_CallbackServer) error {
	ctx := stream.Context()
	implicit := ctx.Value("implicit").([]byte)

	id := ctx.Value("id").(uuid.UUID)
	if id == uuid.Nil {
		var err error
		id, err = uuid.NewRandom()
		if err != nil {
			return grpc.Errorf(
				codes.Internal, "failed to create new UUID: %s",
				err.Error(),
			)
		}

		t, err := s.newToken(id, "/Control/Callback", implicit)
		if err != nil {
			return grpc.Errorf(
				codes.Internal, "failed to issue new token: %s",
				err.Error(),
			)
		}

		err = stream.Send(t)
		if err != nil {
			return fmt.Errorf(
				"failed to send data over the stream: %s",
				err.Error(),
			)
		}
	}

	// TODO: update database
	// ...
	// ...

	close := make(chan struct{})
	go func() {
		p := ping.New(nil)
		p.Start(stream, &pb.Token{})
		close <- struct{}{}
	}()

	// TODO: wait for other new token issue request
	select {
	case <-close:
		// close the stream
		// update the database
	}

	return nil
}

func (s *service) newToken(
	id uuid.UUID,
	resource string,
	implicit []byte) (*pb.Token, error) {

	p := token.Params{
		Issuer:   "/Control",
		Audience: resource,
		Footer:   resource,
		Body: map[string]any{
			"role": token.RoleDaemon,
			"id":   id,
		},
	}
	t, err := token.New(p)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create new token: %s",
			err.Error(),
		)
	}

	enc, err := t.Encrypt(secret, implicit)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to encrypt token: %s",
			err.Error(),
		)
	}

	return &pb.Token{T: enc}, nil
}

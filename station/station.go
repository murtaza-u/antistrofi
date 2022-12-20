package station

import (
	"fmt"
	"net"

	pb "github.com/murtaza-u/antistrofi/proto/gen/go"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type station struct {
	Server *grpc.Server
	ln     net.Listener
}

// New returns a new station from the given parameters.
func New(p Params) (*station, error) {
	err := p.validate()
	if err != nil {
		return nil, err
	}

	secret = p.Secret

	ln, err := newListener(p.Port)
	if err != nil {
		return nil, err
	}

	return &station{
		ln:     *ln,
		Server: newServer(p.Opts),
	}, nil
}

// Start starts the station server. Register all your services before
// calling this method.
func (s *station) Start(reflect bool) error {
	pb.RegisterControlServer(s.Server, &service{})
	if reflect {
		reflection.Register(s.Server)
	}
	return s.Server.Serve(s.ln)
}

func newServer(opt ...grpc.ServerOption) *grpc.Server {
	i := newIntercept()
	srv := grpc.NewServer(grpc.StreamInterceptor(i.stream), opt)
	return srv
}

func newListener(port string) (*net.Listener, error) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to listen on port %s: %s", port, err.Error(),
		)
	}
	return &ln, nil
}

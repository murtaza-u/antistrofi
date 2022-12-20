package station

import (
	"context"

	"github.com/murtaza-u/antistrofi/token"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type serverStream struct {
	grpc.ServerStream
	ctx context.Context
}

type intercept struct{}

func newIntercept() *intercept {
	return &intercept{}
}

func (i intercept) getID(ctx context.Context) (*uuid.UUID, error) {
	implicit := ctx.Value("implicit").([]byte)
	if implicit == nil || len(implicit) == 0 {
		return nil, grpc.Errorf(
			codes.InvalidArgument,
			"missing implicit bytes",
		)
	}

	enc := ctx.Value("token").(string)
	if enc == "" {
		return nil, nil
	}

	t, err := token.Decrypt(enc, secret, implicit)
	if err != nil {
		return nil, grpc.Errorf(
			codes.InvalidArgument,
			"failed to decrypt token: %s", err.Error(),
		)
	}

	id := new(uuid.UUID)
	err = t.Get("id", id)
	if err != nil {
		return nil, grpc.Errorf(
			codes.InvalidArgument,
			err.Error(),
		)
	}

	return id, nil
}

func (i intercept) stream(
	srv any,
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {

	ctx := ss.Context()

	if info.FullMethod == "/Control/Callback" {
		id, err := i.getID(ctx)
		if err != nil {
			return err
		}

		if id != nil {
			ctx = context.WithValue(ctx, "id", id)
		}
	}

	return handler(srv, &serverStream{
		ServerStream: ss,
		ctx:          ctx,
	})
}

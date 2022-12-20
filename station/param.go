package station

import (
	"errors"
	"strconv"

	"google.golang.org/grpc"
)

// DefaultPort is the default port station listens on if nothing is
// provided.
const DefaultPort = "1709"

var (
	ErrMissingServer = errors.New("missing server object in params")
	ErrMissingSecret = errors.New("missing secret key in params")
	ErrInvalidPort   = errors.New("invalid port number")
)

// Params are a list of options to configure the station.
type Params struct {
	Opts    []grpc.ServerOption
	Port    string
	Reflect bool
	Secret  []byte
}

func (p *Params) validatePort() error {
	n, err := strconv.Atoi(p.Port)
	if err != nil {
		return ErrInvalidPort
	}

	if n <= 0 || n > 65535 {
		return ErrInvalidPort
	}

	return nil
}

func (p *Params) validate() error {
	if p.Secret == nil || len(p.Secret) == 0 {
		return ErrMissingServer
	}

	if p.Port == "" {
		p.Port = DefaultPort
	}

	err := p.validatePort()
	if err != nil {
		return err
	}

	return nil
}

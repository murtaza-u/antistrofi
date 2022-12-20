package ping

import "time"

// Stream is an interface that wraps the Send method. All GRPC service
// streams satisfy this interface.
type Stream interface {
	// SendMsg sends a message across the stream
	SendMsg(any) error
}

type ping struct {
	*Cfg
	stop chan struct{}
}

// New instantiates a new ping instance from the provided config.
func New(cfg *Cfg) *ping {
	if cfg == nil {
		cfg = new(Cfg)
	}

	if cfg.Interval == 0 {
		cfg.Interval = DefaultInterval
	}

	return &ping{
		Cfg:  cfg,
		stop: make(chan struct{}),
	}
}

// Stop terminates the pings and forces the *ping.Start method to return
func (p *ping) Stop() {
	p.stop <- struct{}{}
}

// Start pings the stream at regular intervals. It blocks until the
// *ping.Stop method is called or an error occurs while sending pings
// across the stream.
//
// obj is the protocol buffer object the stream can carry over the wire.
func (p *ping) Start(s Stream, obj any) error {
	t := time.NewTicker(p.Interval)
	defer t.Stop()

	for {
		ticked := p.wait(t)
		if !ticked {
			return nil
		}

		err := p.send(s, obj)
		if err != nil {
			return err
		}

		t.Reset(p.Interval)
	}
}

func (ping) send(s Stream, obj any) error {
	return s.SendMsg(obj)
}

func (p *ping) wait(t *time.Ticker) bool {
	select {
	case <-p.stop:
		return false
	case <-t.C:
		return true
	}
}

package store

import (
	"context"
	"errors"

	"go.etcd.io/etcd/clientv3"
)

var (
	ErrWatchCanceled = errors.New("watch canceled by the server")
	ErrWatchClosed   = errors.New("watch closed by the client")
)

// Watcher defines behaviour all watchers must have.
type Watcher interface {
	// Start starts watching specified resource
	Start()
	// Stop stops the watcher
	Stop()
	// EventChan returns events as they occur.
	EventChan() <-chan *Event
	// ErrChan returns an error, if any. The watcher is terminated on
	// receiving an error.
	ErrChan() <-chan error
}

type watch struct {
	wc     clientv3.WatchChan
	ctx    context.Context
	cancel context.CancelFunc
	close  chan struct{}
	evC    chan *Event
	errC   chan error
}

func (w watch) EventChan() <-chan *Event {
	return w.evC
}

func (w watch) ErrChan() <-chan error {
	return w.errC
}

func (w watch) Stop() {
	w.close <- struct{}{}
}

func (w watch) Start() {
	defer w.cancel()

	for {
		select {
		case <-w.close:
			return
		case res := <-w.wc:
			if res.Canceled {
				w.sendErr(ErrWatchCanceled)
				return
			}

			err := res.Err()
			if err != nil {
				w.sendErr(err)
				return
			}

			evs := res.Events
			if len(evs) == 0 {
				w.sendErr(ErrWatchClosed)
				return
			}

			for _, ev := range evs {
				w.sendEvent(parseEvent(ev))
			}
		}
	}
}

func (w watch) sendErr(err error) {
	w.errC <- err
}

func (w watch) sendEvent(ev *Event) {
	w.evC <- ev
}

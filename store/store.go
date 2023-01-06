package store

import (
	"context"
	"errors"
	"time"

	"go.etcd.io/etcd/clientv3"
)

// Timeout defines the default database dial timeout.
const Timeout = time.Second * 30

var ErrKeyNotFound = errors.New("key not found")

// Storer defines behaviour all storage implementations must have.
type Storer interface {
	Put(string, string) error
	Get(string) (string, error)
	Watcher(string) Watcher
}

type store struct {
	client *clientv3.Client
}

// New returns a new storage implememtation.
func New(ends []string) (Storer, error) {
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   ends,
		DialTimeout: Timeout,
	})
	if err != nil {
		return nil, err
	}

	return &store{c}, nil
}

// Put creates/updates a key-value pair.
func (s store) Put(k, v string) error {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	_, err := s.client.Put(ctx, k, v)
	if err != nil {
		return err
	}

	return nil
}

// Get returns the value associated with the key. Returns an error if
// the key does not exist.
func (s store) Get(k string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), Timeout)
	defer cancel()

	res, err := s.client.Get(ctx, k)
	if err != nil {
		return "", err
	}

	if len(res.Kvs) == 0 {
		return "", ErrKeyNotFound
	}

	kv := res.Kvs[0]
	return string(kv.Value), nil
}

// Watcher returns a new watcher instance that can be started to watch
// over the given key.
func (s store) Watcher(k string) Watcher {
	ctx, cancel := context.WithCancel(context.Background())
	wc := s.client.Watch(ctx, k)
	return &watch{
		wc:     wc,
		ctx:    ctx,
		cancel: cancel,
		close:  make(chan struct{}),
		evC:    make(chan *Event),
		errC:   make(chan error),
	}
}

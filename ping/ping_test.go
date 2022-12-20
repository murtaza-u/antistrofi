package ping_test

import (
	"errors"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/murtaza-u/antistrofi/ping"
)

type xstream struct{}
type ystream struct{}

func (xstream) SendMsg(any) error {
	log.Println("received a ping")
	return nil
}

func (ystream) SendMsg(any) error {
	return errors.New("connection closed")
}

func TestGoodPing(t *testing.T) {
	cfg := &ping.Cfg{Interval: time.Millisecond}
	p := ping.New(cfg)

	if p.Interval != time.Millisecond {
		t.Errorf(
			"incorrectly set ping interval. Expected: %v | Got: %v",
			time.Millisecond, p.Interval,
		)
	}

	s := new(xstream)
	var wg sync.WaitGroup
	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		err := p.Start(s, nil)
		if err != nil {
			t.Errorf("Start: %s", err.Error())
		}
		wg.Done()
	}(&wg)

	time.Sleep(time.Millisecond * 3)
	p.Stop()
	wg.Wait()
}

func TestBadPing(t *testing.T) {
	cfg := &ping.Cfg{Interval: time.Millisecond}
	p := ping.New(cfg)

	if p.Interval != time.Millisecond {
		t.Errorf(
			"incorrectly set ping interval. Expected: %v | Got: %v",
			time.Millisecond, p.Interval,
		)
	}

	s := new(ystream)
	err := p.Start(s, nil)
	if err == nil {
		t.Errorf("Start: Expected: an error | Got: nil")
	}
}

/*
Daemon service:
Spawns a number of daemons to listen for requests, which process in different goroutines, given the Routine
  as the worker callback.

Example

	service := NewService(10, MyRoutine, errorNotifier)

	// will block until a daemon receives these messages (no process, just receive).
	service.Send(99)
	service.Send(22)
	service.Send(97)
	service.Send("hello")
	service.Send(42)

	// will block until all daemons exit.
	service.Close()
 */
package backend

import (
	"time"

	"github.com/mijara/statspout/log"
	"errors"
)

type Routine func(interface{}) error
type ErrNotifier func(error)

type Service struct {
	daemons int
	r       Routine
	errNot  ErrNotifier

	closeChan chan bool
	pipe      chan interface{}
}

func daemon(routine Routine, pipe chan interface{}, close chan bool, errNotifier ErrNotifier) {
	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case error:
				errNotifier(t)
			case string:
				errNotifier(errors.New(t))
			}

			// replace the dead daemon with a new one.
			go daemon(routine, pipe, close, errNotifier)

			log.Info.Printf("Daemon died, spawned another one in place.")
		}
	}()

	for {
		select {
		case <-close:
			return
		case req := <-pipe:
			if err := routine(req); err != nil {
				errNotifier(err)
			}
		}
	}
}

func NewService(n int, r Routine, errNot ErrNotifier) *Service {
	closeChan := make(chan bool)
	pipe := make(chan interface{})

	for i := 0; i < n; i++ {
		go daemon(r, pipe, closeChan, errNot)
	}

	log.Info.Printf("%d daemons started.", n)

	return &Service{
		daemons:   n,
		pipe:      pipe,
		closeChan: closeChan,
	}
}

func (s *Service) Send(feed interface{}) {
	s.pipe <- feed
}

func (s *Service) Close() {
	for i := s.daemons; i > 0; i-- {
		s.closeChan <- true
	}

	close(s.pipe)
	close(s.closeChan)

	time.Sleep(time.Millisecond * 500)
}

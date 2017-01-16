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

func daemon(r Routine, pipe chan interface{}, close chan bool, errNotifier ErrNotifier) {
	for {
		select {
		case <-close:
			return
		case req := <-pipe:
			if err := r(req); err != nil {
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

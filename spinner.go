package spinner

import (
	"fmt"
	"os"
	"time"
)

type Spinner struct {
	duration time.Duration
	trigger  interface{}
	sentinel interface{}
	closed   bool
	// sigChan: is the pipe that receives the start and stop
	sigChan chan interface{}
}

func New(freq int64) *Spinner {
	if freq < 1 {
		freq = 10
	}
	sp := Spinner{
		duration: time.Duration(1e9 / freq),
		// sigChan will be created on .Start()
		sigChan:  nil,
		sentinel: nil,
	}

	sp.trigger = &sp
	return &sp
}

func (s *Spinner) Start() error {
	err := s.spin()
	if err == nil {
		s.sigChan <- s.trigger
	}
	return err
}

func (s *Spinner) Stop() {
	if !s.closed && s.sigChan != nil {
		s.sigChan <- s.sentinel
		close(s.sigChan)
		s.closed = true
	}
}

func (s *Spinner) Reset() {
	s.Stop()
	s.sigChan = nil
	s.closed = true
}

func (s *Spinner) Duration() time.Duration {
	return s.duration
}

func (s *Spinner) spin() error {
	if s.sigChan != nil { // Already in use
		return fmt.Errorf("already in use")
	}
	s.sigChan = make(chan interface{})
	go func() {
		// Block till the first symbol comes through
		<-s.sigChan

		var spinner = map[string]string{
			" ⏳ ": "\033[47m",
			" ⌛ ": "\033[48m",
		}
		throttle := time.Tick(s.duration)
		running := true
		for running {
			for segment, color := range spinner {
				select {
				case in := <-s.sigChan:
					if in == s.sentinel {
						running = false
						break
					}
				default:
				}
				// Print it to stderr to avoid symbol getting into piped content
				fmt.Fprintf(os.Stderr, "%s%s\033[00m\r", color, segment)
				<-throttle
			}
		}
	}()
	return nil
}

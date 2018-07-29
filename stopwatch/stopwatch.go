package stopwatch

import (
	"time"
)

// A Stopwatch records how long events take
type Stopwatch struct {
	startTime time.Time
	endTime   time.Time
	duration  time.Duration
}

// New creates a new stopwatch
func New() *Stopwatch {
	return new(Stopwatch)
}

// Start begins the stopwatch
func (s *Stopwatch) Start() {
	s.startTime = time.Now()
}

// Stop finishs the stopwatch
func (s *Stopwatch) Stop() {
	s.endTime = time.Now()
	s.duration = s.endTime.Sub(s.startTime)
}

// Duration returns the elapsed time on the stopwatch
func (s *Stopwatch) Duration() time.Duration {
	return s.duration
}

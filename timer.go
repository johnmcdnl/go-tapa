package tapa

import (
	"time"

	"github.com/google/uuid"
)

type Timer struct {
	ID       string        `json:"-"`
	Start    time.Time     `json:"-"`
	End      time.Time     `json:"-"`
	Duration time.Duration `json:"duration"`
}

func newTimer() *Timer {
	return &Timer{
		ID: uuid.New().String(),
	}
}

func (t *Timer) start() {
	t.Start = time.Now()
}

func (t *Timer) stop() {
	t.End = time.Now()
	t.Duration = t.End.Sub(t.Start)
}

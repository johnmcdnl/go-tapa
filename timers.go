package tapa

import (
	"math"
	"time"
)

type Timers struct {
	Size       int           `json:"size"`
	Timers     []*Timer      `json:"-"`
	Cumulative time.Duration `json:"cumulative"`
	Mean       time.Duration `json:"mean"`
	StdDev     time.Duration `json:"std_dev"`
	Variance   time.Duration `json:"variance"`
	Max        time.Duration `json:"max"`
	Min        time.Duration `json:"min"`
}

func (t *Timers) Add(timer *Timer) {
	t.Timers = append(t.Timers, timer)
	t.Size++
}

func (t *Timers) calculate() {
	t.calcCumulative()
	t.calcMean()
	t.calcVariance()
	t.calcStdDev()
	t.calcMax()
	t.calcMin()
}

func (t *Timers) calcCumulative() time.Duration {
	t.Cumulative = 0
	for _, d := range t.Timers {
		t.Cumulative += d.Duration
	}
	return t.Cumulative
}

func (t *Timers) calcMean() time.Duration {
	if len(t.Timers) == 0 {
		return 0
	}
	t.Mean = t.calcCumulative() / time.Duration(len(t.Timers))
	return t.Mean
}

func (t *Timers) calcVariance() time.Duration {
	if len(t.Timers) == 0 {
		return 0
	}
	mean := t.calcMean()

	var sumValueSquared time.Duration
	for _, d := range t.Timers {
		d1 := d.Duration - mean
		sumValueSquared += d1 * d1
	}
	t.Variance = time.Duration(sumValueSquared / time.Duration(len(t.Timers)))
	return t.Variance
}

func (t *Timers) calcStdDev() {
	t.StdDev = time.Duration(math.Sqrt(float64(t.calcVariance())))
}

func (t *Timers) calcMax() time.Duration {
	for _, d := range t.Timers {
		if d.Duration > t.Max {
			t.Max = d.Duration
		}
	}
	return t.Max
}

func (t *Timers) calcMin() time.Duration {
	t.Min = time.Duration(math.MaxInt64)
	for _, d := range t.Timers {
		if d.Duration < t.Min {
			if d.Duration == 0 {
				panic("why is there an empty duration")
			}
			t.Min = d.Duration
		}
	}
	return t.Min
}

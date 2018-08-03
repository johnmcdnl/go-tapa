package tapa

import (
	"fmt"
	"math"
	"time"
)

// Duration is time.Duration with pretty print for JSON
type Duration time.Duration

// MarshalJSON marshals a duration to JSON in milliseconds
func (d Duration) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%f\"", time.Duration(d).Seconds()*1000)
	return []byte(stamp), nil
}

// Statistic provides data to analyse a test by
type Statistic struct {
	durations  []time.Duration
	Name       string   `json:"name"`
	ErrorCount int      `json:"error_count"`
	Cumulative Duration `json:"cumulative"`
	Mean       Duration `json:"mean"`
	StdDev     Duration `json:"std_dev"`
	Variance   Duration `json:"-"`
	Max        Duration `json:"max"`
	Min        Duration `json:"min"`
}

func (s *Statistic) calculate() {
	if len(s.durations) == 0 {
		return
	}
	s.setCumulative()
	s.setMean()
	s.setVariance()
	s.setStdDev()
	s.setMax()
	s.setMin()
}

func (s *Statistic) setCumulative() {
	for _, d := range s.durations {
		s.Cumulative += Duration(d)
	}
}

func (s *Statistic) setMean() {
	if s.Cumulative == 0 {
		s.setCumulative()
	}
	s.Mean = s.Cumulative / Duration(len(s.durations))
}

func (s *Statistic) setStdDev() {
	if s.Variance == 0 {
		s.setVariance()
	}
	s.StdDev = Duration(math.Sqrt(float64(s.Variance)))
}

func (s *Statistic) setVariance() {
	if s.Mean == 0 {
		s.setMean()
	}
	var sumValueSquared Duration
	for _, d := range s.durations {
		d1 := Duration(d) - s.Mean
		sumValueSquared += d1 * d1
	}
	s.Variance = Duration(sumValueSquared / Duration(len(s.durations)))
}

func (s *Statistic) setMax() {
	s.Max = 0
	for _, d := range s.durations {
		if Duration(d) > s.Max {
			s.Max = Duration(d)
		}
	}
}

func (s *Statistic) setMin() {
	s.Min = Duration(math.MaxInt64)
	for _, d := range s.durations {
		if Duration(d) < s.Min {
			if d == 0 {
				panic("why is there an empty duration")
			}
			s.Min = Duration(d)
		}
	}
}

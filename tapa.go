package tapa

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"fmt"
)

type Tapa struct {
	*Timer                        `json:"timer"`
	*Timers                       `json:"timers"`
	ConcurrentUsers int           `json:"concurrent_users"`
	RequestsPerUser int           `json:"requests_per_user"`
	TotalRequests   int           `json:"total_requests"`
	Request         *http.Request `json:"-"`
}

func New() *Tapa {
	return &Tapa{
		Timer:  newTimer(),
		Timers: new(Timers),
	}

}

func (t *Tapa) String() string {
	j, err := json.MarshalIndent(t, "", "\t")
	if err != nil {
		panic(errors.Wrap(err, "failed to get a tapa.String()"))
	}
	return string(j)
}

func (t *Tapa) Run() {
	defer t.Timer.stop()
	t.Timer.start()

	var count int
	for i := 1; i <= t.ConcurrentUsers; i++ {
		for j := 1; j <= t.RequestsPerUser; j++ {
			count++
			if count%100 == 0 {
				fmt.Println(count, "/", t.ConcurrentUsers*t.RequestsPerUser)
			}
			logrus.Info("User %d Request %d Status: %d \n", i, j)
			t.doRequest()
		}
	}

}

func (t *Tapa) doRequest() {
	t.TotalRequests++
	var timer = newTimer()
	timer.start()
	resp, err := http.DefaultClient.Do(t.Request)
	timer.stop()
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	t.Timers.Add(timer)
}

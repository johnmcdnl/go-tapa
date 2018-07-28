package tapa

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"io"
	"github.com/sirupsen/logrus"
	"time"
)

type Tapa struct {
	*Timer  `json:"timer"`
	*Timers `json:"timers"`
	ErrorCount      int
	concurrentUsers int
	requestsPerUser int
	totalRequests   int

	request    *http.Request
	expectFunc func(r *http.Response) bool
}

func New(users, requests int) *Tapa {
	return &Tapa{
		Timer:           newTimer(),
		Timers:          new(Timers),
		concurrentUsers: users,
		requestsPerUser: requests,
		totalRequests:   users * requests,
	}

}

func (t *Tapa) String() string {
	j, err := json.MarshalIndent(t, "", "\t")
	if err != nil {
		panic(errors.Wrap(err, "failed to get a tapa.String()"))
	}
	return string(j)
}

func (t *Tapa) AddRequest(method, url string, body io.Reader) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err)
	}
	t.request = req
}

func (t *Tapa) AddExpectation(fn func(resp *http.Response) bool) {
	t.expectFunc = fn
}

func (t *Tapa) Run() {
	t.Timer.start()
	t.run()
	t.Timer.stop()
	t.calculate()
}

func (t *Tapa) run() {

	t.warmUp()

	jobs := make(chan *http.Request, t.concurrentUsers)
	results := make(chan *Timer, t.concurrentUsers*t.requestsPerUser)
	for w := 1; w <= t.concurrentUsers; w++ {
		go t.addRequestToQueue(jobs, results)
	}

	for j := 1; j <= t.concurrentUsers*t.requestsPerUser; j++ {
		req := *t.request
		jobs <- &req
	}
	close(jobs)

	for a := 1; a <= t.concurrentUsers*t.requestsPerUser; a++ {
		t.Add(<-results)
	}
}

func (t *Tapa) addRequestToQueue(jobs <-chan *http.Request, results chan<- *Timer) {
	for req := range jobs {
		timer := newTimer()
		timer.start()
		resp, err := t.doRequest(req)
		timer.stop()
		if err != nil {
			t.ErrorCount++
		}
		if !t.expect(resp) {
			t.ErrorCount++
		}

		results <- timer
	}
}

func (t *Tapa) doRequest(req *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(req)
}

func (t *Tapa) warmUp() {
	logrus.Debugln("warmUp() Started")
	jobs := make(chan *http.Request, t.concurrentUsers)
	results := make(chan *Timer, t.concurrentUsers)
	for w := 1; w <= t.concurrentUsers; w++ {
		go t.addRequestToQueue(jobs, results)
	}

	for j := 1; j <= t.concurrentUsers; j++ {
		req := *t.request
		jobs <- &req
	}
	close(jobs)

	for a := 1; a <= t.concurrentUsers; a++ {
		<-results
	}
	time.Sleep(500 * time.Millisecond)

	logrus.Debugln("warmUp() Finished")
}

func (t *Tapa) expect(resp *http.Response) bool {
	return t.expectFunc(resp)
}

package tapa

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/cheggaaa/pb.v1"
)

//Tapa needs a TODO comment
type Tapa struct {
	*Timer          `json:"timer"`
	*Timers         `json:"timers"`
	ErrorCount      int
	concurrentUsers int
	requestsPerUser int
	totalRequests   int
	progressBar     *pb.ProgressBar
	request         *http.Request
	expectFunc      []func(r *http.Response) bool
}

func newProgressBar(total int) *pb.ProgressBar {
	bar := pb.New(total)
	bar.Width = 120
	return bar
}

// New tapa
func New(users, requests int) *Tapa {

	return &Tapa{
		Timer:           newTimer(),
		Timers:          new(Timers),
		concurrentUsers: users,
		requestsPerUser: requests,
		totalRequests:   users * requests,
		progressBar:     newProgressBar(users * requests),
	}

}

func (t *Tapa) reset() {
	t.Timers = new(Timers)
}

func (t *Tapa) String() string {
	j, err := json.MarshalIndent(t, "", "\t")
	if err != nil {
		panic(errors.Wrap(err, "failed to get a tapa.String()"))
	}
	return string(j)
}

// AddRequest adds the HTTP endpoint in test to the suite
func (t *Tapa) AddRequest(method, url string, header http.Header, body io.Reader) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err)
	}
	t.request = req
}

// AddExpectation adds an expectation that shoudl be satisfied by every resposne
func (t *Tapa) AddExpectation(fn func(resp *http.Response) bool) {
	t.expectFunc = append(t.expectFunc, fn)
}

//Run executes the tests
func (t *Tapa) Run() {
	t.Timer.start()
	t.warmUp()
	t.reset()
	t.run()
	t.Timer.stop()
	t.calculate()
}

func (t *Tapa) run() {
	t.progressBar.Start()
	defer t.progressBar.Finish()

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
		t.progressBar.Increment()
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
	origBar := t.progressBar
	defer func() {
		t.progressBar = origBar
	}()

	t.progressBar = newProgressBar(t.concurrentUsers)
	t.progressBar.Start()
	defer t.progressBar.Finish()

	logrus.Debugln("warmUp() Started")
	jobs := make(chan *http.Request, t.concurrentUsers)
	results := make(chan *Timer, t.concurrentUsers)

	for w := 0; w < int(math.Ceil(float64(t.concurrentUsers)/float64(7))); w++ {
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

	logrus.Debugln("warmUp() Finished")
}

func (t *Tapa) expect(resp *http.Response) bool {
	for _, fn := range t.expectFunc {
		if !fn(resp) {
			return false
		}
	}
	return true
}

// Report generates and outputs a report of the statistics
func (t *Tapa) Report() {
	fmt.Println("t.Mean", t.Mean)
	fmt.Println("t.StdDev", t.StdDev)
	fmt.Println("t.Min", t.Min)
	fmt.Println("t.Max", t.Max)
	fmt.Println("t.Duration", t.Duration)
	fmt.Println("t.ErrorCount", t.ErrorCount)
}

package tapa

import (
	"encoding/json"
	"fmt"

	pb "gopkg.in/cheggaaa/pb.v1"

	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/johnmcdnl/go-tapa/stopwatch"
	"github.com/tidwall/gjson"
)

// Tapa is a collection of tests
type Tapa struct {
	stopwatch       stopwatch.Stopwatch
	client          *http.Client
	requests        []*Request
	users           int
	requestsPerUser int
	delayMin        time.Duration
	delayMax        time.Duration
	warmUpPeriod    time.Duration
	rand            *rand.Rand
	concurrency     int
}

// New creates a new suite with default params
func New() *Tapa {
	var t Tapa
	t.concurrency = 2
	t.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	t.WithClient(http.DefaultClient)

	return &t
}

// WithClient allows the consumer to provide their own client
func (t *Tapa) WithClient(client *http.Client) {
	t.client = client
}

// WithConcurrency defines how many tests can run in parallel
func (t *Tapa) WithConcurrency(concurrency int) {
	t.concurrency = concurrency
}

// WithWarmUpPeriod defines how long to wait before starting each test
func (t *Tapa) WithWarmUpPeriod(warmUpPeriod time.Duration) {
	t.warmUpPeriod = warmUpPeriod
}

// WithRequest adds a request to the suite
func (t *Tapa) WithRequest(req *Request) {
	if req.users == 0 {
		req.WithUsers(t.users)
	}
	if req.requestsPerUser == 0 {
		req.WithRequestsPerUser(t.requestsPerUser)
	}
	if req.client == nil {
		req.WithClient(t.client)
	}

	if req.delayMin == 0 {
		req.delayMin = t.delayMin
	}
	if req.delayMax == 0 {
		req.delayMax = t.delayMax
	}

	t.requests = append(t.requests, req)
}

// WithDelay adds a randomised delay between requests
func (t *Tapa) WithDelay(min, max time.Duration) {
	t.delayMin = min
	t.delayMax = max
}

// WithUsers defines how many concurrent users there will be
func (t *Tapa) WithUsers(users int) {
	t.users = users
}

// WithRequestsPerUser gives a default number of requests per user
func (t *Tapa) WithRequestsPerUser(reqsPerUser int) {
	t.requestsPerUser = reqsPerUser
}

func (t *Tapa) getDelay() time.Duration {
	if t.delayMin == t.delayMin {
		return t.delayMin
	}
	return time.Duration(t.rand.Intn(int(t.delayMax-t.delayMin)) + int(t.delayMin))
}

// Run executes the suite
func (t *Tapa) Run() {

	var progressBars []*pb.ProgressBar
	for _, req := range t.requests {
		size := req.users * req.requestsPerUser
		req.progressBar = pb.New(size)
		req.progressBar.Width = 120
		progressBars = append(progressBars, req.progressBar)
	}

	pool := pb.NewPool(progressBars...)
	pool.Start()
	defer pool.Stop()

	jobs := make(chan *Request, t.concurrency)
	responses := make(chan bool, len(t.requests))

	for worker := 1; worker <= t.concurrency; worker++ {
		go t.runRequest(jobs, responses)
	}

	for _, req := range t.requests {
		jobs <- req
	}
	close(jobs)
	for worker := 1; worker <= len(t.requests); worker++ {
		<-responses
	}

}

func (t *Tapa) runRequest(requests chan *Request, finish chan bool) {
	for req := range requests {
		time.Sleep(t.warmUpPeriod)
		req.execute()
		finish <- true
	}
}

// Report outputs stats
func (t *Tapa) Report() {

	var stats []*Statistic
	for _, req := range t.requests {
		var stat Statistic
		stat.Name = req.Name()
		stat.durations = req.durations
		stat.ErrorCount = len(req.durations)
		stat.calculate()
		stats = append(stats, &stat)

	}

	jsonStats := toJSON(stats)

	ioutil.WriteFile("zReport.json", []byte(jsonStats), os.ModePerm)

	fmt.Println(gjson.Parse(jsonStats).Get("#.mean").String())
	fmt.Println(gjson.Parse(jsonStats).Get("#.error_count").String())

}

func toJSON(i interface{}) string {
	j, err := json.MarshalIndent(i, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(j)
}

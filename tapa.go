package tapa

import (
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/johnmcdnl/go-tapa/stopwatch"
	pb "gopkg.in/cheggaaa/pb.v1"
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
	rand            *rand.Rand
}

// New creates a new suite with default params
func New() *Tapa {
	var t Tapa
	t.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	t.WithClient(http.DefaultClient)

	return &t
}

// WithClient allows the consumer to provide their own client
func (t *Tapa) WithClient(client *http.Client) {
	t.client = client
}

// WithRequest adds a request to the suite
func (t *Tapa) WithRequest(req *Request) {
	if req.users == 0 || t.users != 0 {
		req.users = t.users
	}
	if req.requestsPerUser == 0 || t.requestsPerUser != 0 {
		req.requestsPerUser = t.requestsPerUser
	}
	if req.client == nil || t.client != nil {
		req.client = t.client
	}

	if req.delayMin == 0 || t.delayMin != 0 {
		req.delayMin = t.delayMin
	}
	if req.delayMax == 0 || t.delayMax != 0 {
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

	// var pool = pb.NewPool([]*pb.ProgressBar{}...)
	var progressBars []*pb.ProgressBar
	for _, req := range t.requests {
		size := req.users * req.requestsPerUser
		req.progressBar = pb.New(size)
		req.progressBar.Width = 120
		progressBars = append(progressBars, req.progressBar)

	}

	pool := pb.NewPool(progressBars...)

	var wg sync.WaitGroup
	pool.Start()
	for _, req := range t.requests {
		wg.Add(1)
		t.runRequest(&wg, req)
	}
	pool.Stop()
	wg.Wait()
}

func (t *Tapa) runRequest(wg *sync.WaitGroup, req *Request) {
	defer wg.Done()
	req.execute()
}

// Report outputs stats
func (t *Tapa) Report() {
	// for _, req := range t.requests {
	// 	fmt.Println(req.durations, len(req.durations))
	// 	fmt.Println(req.errors, len(req.errors))
	// }
}

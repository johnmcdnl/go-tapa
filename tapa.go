package tapa

import (
	"encoding/json"
	"net/http"

	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Tapa struct {
	*Timer          `json:"timer"`
	*Timers         `json:"timers"`
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
	t.Timer.start()
	t.runConcurrent()
	t.Timer.stop()
	t.calculate()
}

func (t *Tapa) doRequest() error {
	t.TotalRequests++
	var timer = newTimer()
	timer.start()
	resp, err := http.DefaultClient.Do(t.Request)
	timer.stop()
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	t.Timers.Add(timer)
	return nil
}

func (t *Tapa) runConcurrent() {

	var tasks []*Task

	for i := 1; i <= t.ConcurrentUsers; i++ {
		for j := 1; j <= t.RequestsPerUser; j++ {
			tasks = append(tasks, NewTask(t.doRequest))
		}
	}

	p := NewPool(tasks, t.ConcurrentUsers)
	p.Run()
}

type Task struct {
	Err error
	f   func() error
}

func NewTask(f func() error) *Task {
	return &Task{f: f}
}

func (t *Task) Run(wg *sync.WaitGroup) {
	t.Err = t.f()
	wg.Done()
}

type Pool struct {
	Tasks []*Task

	concurrency int
	tasksChan   chan *Task
	wg          sync.WaitGroup
}

func NewPool(tasks []*Task, concurrency int) *Pool {
	return &Pool{
		Tasks:       tasks,
		concurrency: concurrency,
		tasksChan:   make(chan *Task),
	}
}

func (p *Pool) HasErrors() bool {
	for _, task := range p.Tasks {
		if task.Err != nil {
			return true
		}
	}
	return false
}

func (p *Pool) Run() {
	logrus.Debugln("Running %v task(s) at concurrency %v.", len(p.Tasks), p.concurrency)

	for i := 0; i < p.concurrency; i++ {
		go p.work()
	}

	p.wg.Add(len(p.Tasks))
	for _, task := range p.Tasks {
		p.tasksChan <- task
	}

	close(p.tasksChan)

	p.wg.Wait()
}

func (p *Pool) work() {
	for task := range p.tasksChan {
		task.Run(&p.wg)
	}
}

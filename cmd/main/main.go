package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/johnmcdnl/go-tapa"
	"github.com/johnmcdnl/go-tapa/cmd/server"
	"github.com/sirupsen/logrus"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	logrus.SetLevel(logrus.InfoLevel)
	go server.Server(8532)
	time.Sleep(1 * time.Second)
}

func main() {
	var t = tapa.New()
	t.WithUsers(25)
	t.WithRequestsPerUser(1000)
	t.WithConcurrency(1)
	t.WithWarmUpPeriod(2 * time.Second)
	t.WithDelay(time.Duration(3*time.Millisecond), time.Duration(5*time.Millisecond))
	//t.WithDelay(time.Duration(1*time.Millisecond), time.Duration(2*time.Millisecond))
	// t.WithDelay(time.Duration(500*time.Millisecond), time.Duration(500*time.Millisecond))

	for i := 0; i < 3; i++ {
		addRequest(t)
	}

	t.Run()
	t.Report()
}

func addRequest(t *tapa.Tapa) {
	req := tapa.NewRequestMust(
		http.NewRequest(http.MethodGet, "http://localhost:8532", nil),
	)
	req.WithExpectation(func(resp *http.Response) error {
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("expected to get a 200 status code but got %d", resp.StatusCode)
		}
		return nil
	})
	req.WithExpectation(func(resp *http.Response) error {
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("expected to get a 200 status code but got %d", resp.StatusCode)
		}
		return nil
	})

	t.WithRequest(req)
}

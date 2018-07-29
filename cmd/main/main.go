package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/johnmcdnl/go-tapa"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	var t = tapa.New()
	t.WithUsers(50)
	t.WithRequestsPerUser(10)
	t.WithDelay(time.Duration(500*time.Millisecond), time.Duration(500*time.Millisecond))

	req := tapa.NewRequestMust(
		http.NewRequest(http.MethodGet, "http://localhost:8532", nil),
	)
	req.WithExpectation(func(resp *http.Response) error {
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("expected to get a 200 status code but got %d", resp.StatusCode)
		}
		return nil
	})
	t.WithRequest(req)

	req1 := tapa.NewRequestMust(
		http.NewRequest(http.MethodGet, "http://localhost:8532", nil),
	)
	req1.WithExpectation(func(resp *http.Response) error {
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("expected to get a 200 status code but got %d", resp.StatusCode)
		}
		return nil
	})
	t.WithRequest(req1)

	t.Run()
	t.Report()
}

package main

import (
	"io/ioutil"
	"net/http"

	"github.com/johnmcdnl/go-tapa"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	var t = tapa.New(50, 10)
	t.AddRequest(http.MethodGet, "http://localhost:8532", nil, nil)
	t.AddExpectation(func(resp *http.Response) bool {
		if resp.StatusCode != http.StatusOK {
			return false
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		if len(b) != 0 {
			return false
		}
		return true
	})
	t.AddExpectation(func(resp *http.Response) bool {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		if len(b) != 0 {
			return false
		}
		return true
	})
	t.AddExpectation(func(resp *http.Response) bool {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		if len(b) != 0 {
			return false
		}
		return true
	})
	t.Run()

	t.Report()

}

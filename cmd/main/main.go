package main

import (
	"fmt"
	"net/http"

	"github.com/johnmcdnl/go-tapa"
	"github.com/sirupsen/logrus"
	"time"
	"io/ioutil"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	go testEndpoint()
}

func testEndpoint() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
		logrus.Debugln("Hi there, I love %s!", r.URL.Path[1:])
	}

	http.HandleFunc("/", handler)
	logrus.Fatal(http.ListenAndServe(":8532", nil))
}

func main() {
	var t = tapa.New(50, 10)
	t.AddRequest(http.MethodGet, "http://localhost:8532", nil)
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
	t.Run()

	fmt.Println(t)
	fmt.Println("mean", t.Mean)
	fmt.Println(t.StdDev)
	fmt.Println(t.Min)
	fmt.Println(t.Max)
	fmt.Println("suite", t.Duration)
	fmt.Println("errCount", t.ErrorCount)
}

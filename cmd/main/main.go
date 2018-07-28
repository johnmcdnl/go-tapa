package main

import (
	"fmt"
	"net/http"

	"github.com/johnmcdnl/go-tapa"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

func init() {
	logrus.SetLevel(logrus.InfoLevel)
	go testEndpoint()
}

func testEndpoint() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		//time.Sleep(1 * time.Second)
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
	fmt.Println("t.Mean", t.Mean)
	fmt.Println("t.StdDev", t.StdDev)
	fmt.Println("t.Min", t.Min)
	fmt.Println("t.Max", t.Max)
	fmt.Println("t.Duration", t.Duration)
	fmt.Println("t.ErrorCount", t.ErrorCount)
}

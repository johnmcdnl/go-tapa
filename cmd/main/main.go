package main

import (
	"fmt"
	"net/http"

	"github.com/johnmcdnl/go-tapa"
	"github.com/sirupsen/logrus"
	"time"
)

func init() {
	logrus.SetLevel(logrus.InfoLevel)
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
	var t = tapa.New(50, 50)
	t.AddRequest(http.MethodGet, "http://localhost:8532", nil)
	t.Run()

	fmt.Println(t)
	fmt.Println("mean", t.Mean)
	fmt.Println(t.StdDev)
	fmt.Println(t.Min)
	fmt.Println(t.Max)
	fmt.Println("suite", t.Duration)
}

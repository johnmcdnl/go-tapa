package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/johnmcdnl/go-tapa"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.ErrorLevel)
	go testEndpoint()
}

func testEndpoint() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("testEndpoint()")
		fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8532", nil))
}

func main() {
	var t = tapa.New()
	t.ConcurrentUsers = 100
	t.RequestsPerUser = 100
	t.Request = func() *http.Request {
		req, err := http.NewRequest(
			http.MethodGet,
			"http://localhost:8532",
			nil,
		)
		if err != nil {
			panic(err)
		}
		return req
	}()

	t.Run()

	fmt.Println(t)
	fmt.Println(t.Mean)

}

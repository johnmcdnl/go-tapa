package main

import (
	"net/http"
	"github.com/sirupsen/logrus"
)

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		//time.Sleep(1 * time.Second)
		logrus.Debugln("Hi there, I love %s!", r.URL.Path[1:])
	}

	http.HandleFunc("/", handler)
	logrus.Fatal(http.ListenAndServe(":8532", nil))
}

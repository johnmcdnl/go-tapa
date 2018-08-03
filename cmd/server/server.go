package server

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

// Server runs a basic test server on port
func Server(port int) {
	logrus.Infoln("Setting up server on port", port)
	handler := func(w http.ResponseWriter, r *http.Request) {
		logrus.Debugln("Hi there, I love %s!", r.URL)
	}

	http.HandleFunc("/", handler)
	logrus.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

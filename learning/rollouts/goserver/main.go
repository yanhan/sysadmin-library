package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	_port          = 3333
	_apiVersion    = "v1"
	_versionHeader = "version"
)

var (
	_failHalfTheTime     = false
	_failHalfTheTimeLock sync.RWMutex

	_mainRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "myapp_main_requests",
			Help: "Number of requests to the / endpoint, grouped by success or failure",
		},
		[]string{"success"},
	)
	_mainRequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_main_requests_total",
		Help: "Total number of requests to the / endpoint",
	})
)

func requestShouldFail() bool {
	return rand.Intn(100) < 50
}

func main() {
	fmt.Println("vim-go")
	rand.Seed(time.Now().UnixNano())

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_failHalfTheTimeLock.RLock()
		defer _failHalfTheTimeLock.RUnlock()
		_mainRequestsTotal.Inc()
		w.Header().Add(_versionHeader, _apiVersion)
		if _failHalfTheTime && requestShouldFail() {
			_mainRequests.WithLabelValues("false").Inc()
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "An error occurred")
			return
		}
		_mainRequests.WithLabelValues("true").Inc()
		fmt.Fprintf(w, "Hello")
	})
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(_versionHeader, _apiVersion)
		fmt.Fprintf(w, "OK")
	})
	http.HandleFunc("/fail/enable", func(w http.ResponseWriter, r *http.Request) {
		_failHalfTheTimeLock.Lock()
		defer _failHalfTheTimeLock.Unlock()
		_failHalfTheTime = true
		w.Header().Add(_versionHeader, _apiVersion)
		fmt.Fprint(w, "Enabled")
	})
	http.HandleFunc("/fail/status", func(w http.ResponseWriter, r *http.Request) {
		_failHalfTheTimeLock.RLock()
		defer _failHalfTheTimeLock.RUnlock()
		w.Header().Add(_versionHeader, _apiVersion)
		fmt.Fprintf(w, "%t", _failHalfTheTime)
	})
	fmt.Printf("Serving at port %d\n", _port)
	http.ListenAndServe(fmt.Sprintf(":%d", _port), nil)
}

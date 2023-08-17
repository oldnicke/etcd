package httpproxy

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestsIncoming = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "etcd",
			Subsystem: "proxy",
			Name:      "requests_total",
			Help:      "Counter requests incoming by method.",
		}, []string{"method"})

	requestsHandled = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "etcd",
			Subsystem: "proxy",
			Name:      "handled_total",
			Help:      "Counter of requests fully handled (by authoratitave servers)",
		}, []string{"method", "code"})

	requestsDropped = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "etcd",
			Subsystem: "proxy",
			Name:      "dropped_total",
			Help:      "Counter of requests dropped on the proxy.",
		}, []string{"method", "proxying_error"})

	requestsHandlingSec = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "etcd",
			Subsystem: "proxy",
			Name:      "handling_duration_seconds",
			Help:      "Bucketed histogram of handling time of successful events (non-watches), by method (GET/PUT etc.).",

			// lowest bucket start of upper bound 0.0005 sec (0.5 ms) with factor 2
			// highest bucket start of 0.0005 sec * 2^12 == 2.048 sec
			Buckets: prometheus.ExponentialBuckets(0.0005, 2, 13),
		}, []string{"method"})
)

type forwardingError string

const (
	zeroEndpoints         forwardingError = "zero_endpoints"
	failedSendingRequest  forwardingError = "failed_sending_request"
	failedGettingResponse forwardingError = "failed_getting_response"
)

func init() {
	prometheus.MustRegister(requestsIncoming)
	prometheus.MustRegister(requestsHandled)
	prometheus.MustRegister(requestsDropped)
	prometheus.MustRegister(requestsHandlingSec)
}

func reportIncomingRequest(request *http.Request) {
	requestsIncoming.WithLabelValues(request.Method).Inc()
}

func reportRequestHandled(request *http.Request, response *http.Response, startTime time.Time) {
	method := request.Method
	requestsHandled.WithLabelValues(method, strconv.Itoa(response.StatusCode)).Inc()
	requestsHandlingSec.WithLabelValues(method).Observe(time.Since(startTime).Seconds())
}

func reportRequestDropped(request *http.Request, err forwardingError) {
	requestsDropped.WithLabelValues(request.Method, string(err)).Inc()
}

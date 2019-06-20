package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	inFlightGaugeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "http_requests_in_flight_total",
		Help: "Contains the number of requests that are currently in-flight.",
	}, []string{"handler"})

	requestCounterVec = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Counts the total number of HTTP requests against a handler.",
	}, []string{"handler", "code", "method"})

	requestDurationVec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Measures the request durations by handler.",
		Buckets: prometheus.DefBuckets,
	}, []string{"handler", "code", "method"})
)

func init() {
	prometheus.MustRegister(
		inFlightGaugeVec,
		requestCounterVec,
		requestDurationVec)
}

func instrumentHandler(name string, handler http.Handler) http.Handler {
	nameLabel := prometheus.Labels{
		"handler": name,
	}
	return promhttp.InstrumentHandlerInFlight(inFlightGaugeVec.With(nameLabel),
		promhttp.InstrumentHandlerCounter(requestCounterVec.MustCurryWith(nameLabel),
			promhttp.InstrumentHandlerDuration(requestDurationVec.MustCurryWith(nameLabel),
				handler)))
}

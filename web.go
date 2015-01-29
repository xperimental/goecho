package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sort"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

func echoHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "URL: %s\n", req.URL)
	fmt.Fprintln(w, "Header:")

	headers := make([]string, 0, len(req.Header))
	for key := range req.Header {
		headers = append(headers, key)
	}
	sort.Strings(headers)

	for _, key := range headers {
		value := req.Header[key]
		values := `"` + strings.Join(value, `"; "`) + `"`
		fmt.Fprintf(w, "%s -> %s\n", key, values)
	}
}

func versionHandler(version string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; chatset=utf-8")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(struct {
			Version string `json:"version"`
			Runtime string `json:"runtime"`
		}{
			Version: version,
			Runtime: runtime.Version(),
		})

		if err != nil {
			http.Error(w, fmt.Sprintf("error creating JSON: %s", err), http.StatusInternalServerError)
		}
	}
}

func createServer(addr, version string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", prometheus.InstrumentHandlerFunc("echo", echoHandler))
	mux.Handle("/metrics", prometheus.UninstrumentedHandler())
	mux.HandleFunc("/version", prometheus.InstrumentHandlerFunc("version", versionHandler(version)))

	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}

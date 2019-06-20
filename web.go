package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sort"
	"strings"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func echoHandler(hostname string, env []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		showEnv := len(r.URL.Query().Get("env")) > 0

		fmt.Fprintf(w, "URL: %s\n", r.URL)
		fmt.Fprintln(w, "Header:")

		headers := make([]string, 0, len(r.Header))
		for key := range r.Header {
			headers = append(headers, key)
		}
		sort.Strings(headers)

		for _, key := range headers {
			value := r.Header[key]
			values := `"` + strings.Join(value, `"; "`) + `"`
			fmt.Fprintf(w, "%s -> %s\n", key, values)
		}

		if len(hostname) > 0 {
			fmt.Fprintf(w, "\nServer: %s\n", hostname)
		}

		if showEnv && len(env) > 0 {
			fmt.Fprintln(w, "\nEnvironment:")
			for _, e := range env {
				fmt.Fprintln(w, e)
			}
		}
	})
}

func versionHandler(version string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	})
}

func createServer(addr, version, hostname string, env []string) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/", instrumentHandler("echo", echoHandler(hostname, env)))
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/version", instrumentHandler("version", versionHandler(version)))

	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}

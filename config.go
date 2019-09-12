package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	envAddr          = "LISTEN_ADDR"
	envAllowEnv      = "ALLOW_ENV"
	envTLSCert       = "TLS_CERT"
	envTLSKey        = "TLS_KEY"
	envGracefulDelay = "GRACEFUL_DELAY"
)

type tlsConfig struct {
	CertFile string
	KeyFile  string
}

type config struct {
	Addr          string
	AllowEnv      bool
	GracefulDelay time.Duration
	TLS           tlsConfig
}

func parseConfig() (config, error) {
	cfg := config{
		Addr:          ":8080",
		AllowEnv:      false,
		GracefulDelay: 2 * time.Second,
	}
	flag.StringVar(&cfg.Addr, "addr", cfg.Addr, "Address and port to listen on.")
	flag.StringVar(&cfg.TLS.CertFile, "tls-cert", cfg.TLS.CertFile, "Path to TLS certificate file.")
	flag.StringVar(&cfg.TLS.KeyFile, "tls-key", cfg.TLS.KeyFile, "Path to TLS key file.")
	flag.DurationVar(&cfg.GracefulDelay, "graceful-delay", cfg.GracefulDelay, "Delay between receiving a shutdown signal and starting shutdown.")
	flag.BoolVar(&cfg.AllowEnv, "allow-env", cfg.AllowEnv, "Allow retrieval of environment variables.")
	flag.Parse()

	if addr, ok := os.LookupEnv(envAddr); ok {
		cfg.Addr = addr
	}

	if allowEnvRaw, ok := os.LookupEnv(envAllowEnv); ok {
		allowEnv, err := strconv.ParseBool(allowEnvRaw)
		if err != nil {
			return config{}, fmt.Errorf("error parsing %q: %s", envAllowEnv, err)
		}
		cfg.AllowEnv = allowEnv
	}

	if certFile, ok := os.LookupEnv(envTLSCert); ok {
		cfg.TLS.CertFile = certFile
	}

	if keyFile, ok := os.LookupEnv(envTLSKey); ok {
		cfg.TLS.KeyFile = keyFile
	}

	if gracefulDelayRaw, ok := os.LookupEnv(envGracefulDelay); ok {
		gracefulDelay, err := time.ParseDuration(gracefulDelayRaw)
		if err != nil {
			return config{}, fmt.Errorf("error parsing %q: %s", envGracefulDelay, err)
		}
		cfg.GracefulDelay = gracefulDelay
	}

	return cfg, nil
}

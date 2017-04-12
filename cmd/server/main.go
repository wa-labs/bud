package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/wa-labs/bud/pkg/service"
)

// FLAG DOMAIN
// Only func main has the right to decide which flags are available to the user.
// If an environment variable is needed then also make it
// overwritable by a flag.
// https://peter.bourgon.org/go-best-practices-2016/#top-tip-5
const (
	// Transport domain.
	defaultDebugAddr = ":6969"
	defaultHTTPAddr  = ":4200"
	// Environment Domain
	defaultTag      = "Develop"
	defaultLogLevel = "Debug"
)

func main() {

	var (
		// Transport domain.
		debugAddrEnv = service.EnvString("DEBUG_ADDR", defaultDebugAddr)
		httpAddrEnv  = service.EnvString("HTTP_ADDR", defaultHTTPAddr)
		// Environment Domain
		envTag      = service.EnvString("VERSION_TAG", defaultTag)
		envLogLevel = service.EnvString("LOG_LEVEL", defaultLogLevel)
	)

	var (
		// Transport domain.
		debugAddr = flag.String(
			"debugAddr", debugAddrEnv, "Debug, health and metrics listen address",
		)
		httpAddr = flag.String(
			"httpAddr", httpAddrEnv, "HTTP listen address",
		)
		// Environment Domain
		tag = flag.String(
			"tag", envTag, "Build Version Tag",
		)
		logLevel = flag.String(
			"logLevel", envLogLevel, "Debug or Prod log level",
		)
	)
	flag.Parse()

	// MECHANICAL DOMAIN
	errc := make(chan error)
	exitCode := 0
	defer func() {
		os.Exit(exitCode)
	}()

	// ENVIRONMENT DOMAIN
	// INITIALIZE DEPENDENCIES
	// Make dependencies explicit!
	// https://peter.bourgon.org/go-best-practices-2016/#top-tip-9

	logger := service.NewJSONLogger(*logLevel)
	logger.With("caller", service.DefaultCaller)

	// DECLARE SERVICES

	// TRANSPORT DOMAIN
	// HTTP TRANSPORT
	httpMux := chi.NewRouter()
	httpMux.Use(middleware.RequestID)
	httpMux.Use(middleware.RealIP)
	httpMux.Use(middleware.Recoverer)

	httpMux.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Hello, World!"))
	})

	go func() {
		logger.Prod("message", "listening on "+*httpAddr+" (HTTP)")
		errc <- http.ListenAndServe(*httpAddr, httpMux)
	}()

	// DEBUG TRANSPORT
	debugMux := chi.NewRouter()
	debugMux.Mount("/debug", middleware.Profiler())
	debugMux.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Alive and healthy. Version: " + *tag))
	})
	debugMux.Handle("/metrics", promhttp.Handler())
	go func() {
		logger.Prod("message", "listening on "+*debugAddr+" (debug)")
		errc <- http.ListenAndServe(*debugAddr, debugMux)
	}()

	// MECHANICAL DOMAIN
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		logger.Prod(
			"signal", fmt.Sprintf("%s", <-c),
			"msg", "gracefully shutting down",
		)
		errc <- nil
	}()

	if err := <-errc; err != nil {
		logger.Prod("error", err)
		exitCode = 1
	}

}

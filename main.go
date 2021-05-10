package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/mkatychev/go-httpbin/httpbin"
)

const defaultHost = "0.0.0.0"
const defaultPort = 8080

var (
	host          string
	port          int
	maxBodySize   int64
	maxDuration   time.Duration
	httpsCertFile string
	httpsKeyFile  string
)

// main implements the go-httpbin CLI's main() function in a reusable way
func main() {
	host = defaultHost
	port = defaultPort
	httpsCertFile = ""
	httpsKeyFile = ""
	maxBodySize = httpbin.DefaultMaxBodySize
	maxDuration = httpbin.DefaultMaxDuration

	logger := log.New(os.Stderr, "", 0)

	// A hacky log helper function to ensure that shutdown messages are
	// formatted the same as other messages.  See StdLogObserver in
	// httpbin/middleware.go for the format we're matching here.
	serverLog := func(msg string, args ...interface{}) {
		const (
			logFmt  = "time=%q msg=%q"
			dateFmt = "2006-01-02T15:04:05.9999"
		)
		logger.Printf(logFmt, time.Now().Format(dateFmt), fmt.Sprintf(msg, args...))
	}

	h := httpbin.New(
		httpbin.WithMaxBodySize(maxBodySize),
		httpbin.WithMaxDuration(maxDuration),
		httpbin.WithObserver(httpbin.StdLogObserver(logger)),
	)

	listenAddr := net.JoinHostPort(host, strconv.Itoa(port))

	server := &http.Server{
		Addr:    listenAddr,
		Handler: h.Handler(httpbin.HandleFunc("/example", unsortedStrings, "GET")),
	}

	// shutdownCh triggers graceful shutdown on SIGINT or SIGTERM
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	// exitCh will be closed when it is safe to exit, after graceful shutdown
	exitCh := make(chan struct{})

	go func() {
		sig := <-shutdownCh
		serverLog("shutdown started by signal: %s", sig)

		shutdownTimeout := maxDuration + 1*time.Second
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			serverLog("shutdown error: %s", err)
		}

		close(exitCh)
	}()

	var listenErr error
	serverLog("go-httpbin listening on http://%s", listenAddr)
	listenErr = server.ListenAndServe()

	if listenErr != nil && listenErr != http.ErrServerClosed {
		logger.Fatalf("failed to listen: %s", listenErr)
	}

	<-exitCh
	serverLog("shutdown finished")
}

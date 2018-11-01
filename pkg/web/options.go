package web

import (
	"strings"

	"github.com/gorilla/mux"
)

// Options contain the configuration for the web server
type Options struct {
	Port             string
	RequestPerSecond int64
	Burst            int
	Wait             mux.MiddlewareFunc
}

// Option custom type
type Option func(*Options)

// Port defines on which TCP/IP port the web server will run on
func Port(port string) Option {
	if !strings.Contains(port, ":") {
		port = ":" + port
	}

	return func(args *Options) {
		args.Port = port
	}
}

// RequestsPerSecond defines how many HTTP requests can be handled per second of real time.
func RequestsPerSecond(rps int64) Option {
	return func(args *Options) {
		args.RequestPerSecond = rps
	}
}

// Wait specifies we should use the wait protocol versuses the allow.
func Wait() Option {
	return func(args *Options) {
		args.Wait = waitLimitsMW
	}
}

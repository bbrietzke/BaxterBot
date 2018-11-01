package web

import "strings"

// Options contain the configuration for the web server
type Options struct {
	Port             string
	RequestPerSecond int64
	Burst            int
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

// BurstLimit is the maximuim number of incoming requests
func BurstLimit(burst int) Option {
	return func(args *Options) {
		args.Burst = burst
	}
}

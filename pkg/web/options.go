package web

import "strings"

// Options contain the configuration for the web server
type Options struct {
	Port string
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

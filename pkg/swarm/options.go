package swarm

import (
	"os"
	"strings"
)

// Options contain the configuration for the swarm
type Options struct {
	Port       string
	SingleNode bool
	Join       string
	Name       string
	Named      bool
	DataDir    string
	HTTP       string
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

// HTTP defines a custom HTTP port
func HTTP(port string) Option {
	if !strings.Contains(port, ":") {
		port = ":" + port
	}

	return func(args *Options) {
		args.HTTP = port
	}
}

// Join specifies the host and port that we should talk to in order to join the swarm
func Join(host string) Option {
	return func(args *Options) {
		args.SingleNode = false
		args.Join = host
	}
}

// Name specifies the name with which to represent ourselves in the swarm
func Name(name string) Option {
	return func(args *Options) {
		args.Name = name
		args.Named = true
	}
}

// DataDir sets the path to store data
func DataDir(path string) Option {
	if s, err := os.Stat(path); err == nil && s.IsDir() {
		return func(args *Options) {
			args.DataDir = path
		}
	}

	return func(args *Options) {
	}
}

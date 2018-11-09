package store

import (
	"strings"
)

// Argument is the functional implementation of the specification or parameter
type Argument func(*Arguments)

// Arguments are the specifications or parameters to get the storage system up and running.
type Arguments struct {
	port          string
	join          []string
	name          string
	dataDirectory string
	bootstrap     bool
}

// Bootstrap will initialize the replication process.
func Bootstrap() Argument {
	return func(args *Arguments) {
		args.bootstrap = true
	}
}

// Port defines the tcp/ip socket to host the replication
func Port(port string) Argument {
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	return func(args *Arguments) {
		args.port = port
	}
}

// Join represents a single host ( ip address and port ) that we should attempt to contact in order to join a replication network.
func Join(address string) Argument {
	return func(args *Arguments) {
		args.join = append(args.join, address)
	}
}

// Name is the optional name to assign to this node.  If a name is not provided, then one is automatically generated.
func Name(name string) Argument {
	return func(args *Arguments) {
		args.name = name
	}
}

// DataDirectory is the location that we will save the data files, along with replication logs and snapshots.
// If one is not provided, then the values will be saved and in memory and will be lost once the service is turned off.
func DataDirectory(name string) Argument {
	return func(args *Arguments) {
		args.dataDirectory = name
	}
}

func (a *Arguments) persist() bool {
	return len(a.dataDirectory) > 0
}

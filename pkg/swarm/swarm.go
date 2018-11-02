package swarm

import (
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/raft"
)

var (
	swarmer *raft.Raft
	logger  *log.Logger
)

func init() {
	logger = log.New(os.Stdout, "SWARM ", log.LstdFlags|log.Lshortfile)
}

// Start gets the swarm up and running
func Start(options ...Option) error {
	args := &Options{Port: ":21000", SingleNode: true, Join: "", Name: strings.ToUpper(NewName().Haikunate())}
	config := raft.DefaultConfig()
	config.Logger = logger

	for _, o := range options {
		o(args)
	}

	addr, err := net.ResolveTCPAddr("tcp", args.Port)

	if err != nil {
		return err
	}
	config.LocalID = raft.ServerID(args.Name)

	transport, err := raft.NewTCPTransportWithLogger(args.Port, addr, 10, 20*time.Second, logger)

	if err != nil {
		return err
	}

	swarmer, err = raft.NewRaft(config, newStateMachine(), raft.NewInmemStore(), raft.NewInmemStore(), raft.NewInmemSnapshotStore(), transport)

	if err != nil {
		return err
	}

	f := swarmer.BootstrapCluster(raft.Configuration{Servers: []raft.Server{{ID: config.LocalID, Address: transport.LocalAddr()}}})
	return f.Error()
}

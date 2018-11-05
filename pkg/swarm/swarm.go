package swarm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/raft"
)

var (
	swarmer *raft.Raft
	logger  *log.Logger
	myAddr  raft.ServerAddress
	leader  bool
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

	addr, err := net.ResolveTCPAddr("tcp", outboundIP(args.Port))

	if err != nil {
		return err
	}
	config.LocalID = raft.ServerID(args.Name)

	transport, err := raft.NewTCPTransportWithLogger(args.Port, addr, 10, 20*time.Second, logger)
	myAddr = transport.LocalAddr()

	if err != nil {
		return err
	}

	swarmer, err = raft.NewRaft(config, newStateMachine(), raft.NewInmemStore(), raft.NewInmemStore(), raft.NewInmemSnapshotStore(), transport)

	if err != nil {
		return err
	}
	go listenForLeadership(swarmer.LeaderCh())

	if args.SingleNode {
		f := swarmer.BootstrapCluster(raft.Configuration{Servers: []raft.Server{{ID: config.LocalID, Address: myAddr}}})
		return f.Error()
	} else {
		logger.Println(args.Name, outboundIP(args.Port))
		v, err := json.Marshal(Registration{Name: args.Name, Address: outboundIP(args.Port)})
		if err != nil {
			logger.Panicln(err)
		}
		if r, err := http.Post(fmt.Sprintf("http://%s/join", args.Join), "application-type/json", bytes.NewReader(v)); err == nil {
			defer r.Body.Close()
		} else {
			logger.Panicln(err)
		}
	}

	return nil
}

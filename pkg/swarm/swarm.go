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

const (
	applyTimeout             time.Duration = 2 * time.Second
	monitorLeadershipTimeout time.Duration = 10 * time.Minute
)

var (
	swarmer       *raft.Raft
	logger        *log.Logger
	myAddr        raft.ServerAddress
	leader        bool
	grpcPort      string
	httpPort      string
	leaderNetAddr string
)

func init() {
	logger = log.New(os.Stdout, "SWARM ", log.LstdFlags|log.Lshortfile)
	httpPort = ":8080"
	grpcPort = ":8100"
}

// Start gets the swarm up and running
func Start(options ...Option) error {
	args := &Options{Port: ":21000", SingleNode: true, Join: "", Name: strings.ToUpper(NewName().Haikunate()), HTTP: httpPort}
	config := raft.DefaultConfig()
	config.Logger = logger

	for _, o := range options {
		o(args)
	}

	httpPort = args.HTTP
	addr, err := net.ResolveTCPAddr("tcp", outboundIP(args.Port))

	if err != nil {
		return err
	}
	config.LocalID = raft.ServerID(args.Name)

	transport, err := raft.NewTCPTransportWithLogger(outboundIP(args.Port), addr, 10, 20*time.Second, logger)
	if err != nil {
		return err
	}
	myAddr = transport.LocalAddr()

	if err != nil {
		return err
	}

	store, err := raft.NewFileSnapshotStoreWithLogger("/tmp", 2, logger)
	swarmer, err = raft.NewRaft(config, newStateMachine(), raft.NewInmemStore(), raft.NewInmemStore(), store, transport)

	if err != nil {
		return err
	}

	go listenForLeadership(swarmer.LeaderCh())

	if args.SingleNode {
		return swarmer.BootstrapCluster(raft.Configuration{Servers: []raft.Server{{ID: config.LocalID, Address: myAddr}}}).Error()
	}

	logger.Println(args.Name, outboundIP(args.Port))
	v, err := json.Marshal(registration{Name: args.Name, Address: outboundIP(args.Port)})
	if err != nil {
		logger.Panicln(err)
	}
	if r, err := http.Post(fmt.Sprintf("http://%s/join", args.Join), "application-type/json", bytes.NewReader(v)); err == nil {
		logger.Println(r.StatusCode)
		defer r.Body.Close()
	} else {
		logger.Panicln(err)
	}

	return nil
}

// IsLeader reveals if this host is the leader for the swarm.
func IsLeader() bool {
	return leader
}

// LeaderAddr returns the IP address or DNS name of the current cluster leader.
func LeaderAddr() string {
	return leaderNetAddr
}

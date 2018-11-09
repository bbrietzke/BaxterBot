package store

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/hashicorp/golang-lru"

	"github.com/hashicorp/raft"
	"github.com/pkg/errors"
)

const (
	// DefaultReplPort is the default port that replication will happen on
	DefaultReplPort string        = ":8012"
	defaultTimeOut  time.Duration = 3 * time.Second
)

var (
	notImplemented  error
	logger          *log.Logger
	repl            *raft.Raft
	isCurrentLeader bool
	cache           *lru.Cache
)

func init() {
	logger = log.New(os.Stdout, "STORE ", log.LstdFlags)
	notImplemented = errors.New("Not Implemented")
	isCurrentLeader = false
}

// Start creates or returns the existing replicated store with the included parameters.
func Start(arguments ...Argument) error {
	args := &Arguments{port: DefaultReplPort, name: "NA", join: make([]string, 0), dataDirectory: ""}
	cache, _ = lru.NewWithEvict(300, evicted)

	for _, a := range arguments {
		a(args)
	}

	addr, err := net.ResolveTCPAddr("tcp", outboundIP(args.port))
	if err != nil {
		return errors.Wrap(err, "resolving tcp addr")
	}

	transport, err := raft.NewTCPTransportWithLogger(outboundIP(args.port), addr, 10, 20*time.Second, logger)
	if err != nil {
		return errors.Wrap(err, "tcp transport with logger")
	}

	repl, err = raft.NewRaft(raftConfig(args), newStateMachine(args), logStore(args), stableStore(args), snapStore(args), transport)
	if err != nil {
		return errors.Wrap(err, "new raft")
	}

	go leadershipChanges(repl.LeaderCh())
	go bootstrapOrJoin(args, transport.LocalAddr())

	return nil
}

// Stop shutdowns the current store.
func Stop() error {
	return notImplemented
}

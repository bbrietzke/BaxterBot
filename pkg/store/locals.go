package store

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

func outboundIP(port string) string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String() + port
}

func raftConfig(args *Arguments) *raft.Config {
	config := raft.DefaultConfig()
	config.Logger = log.New(os.Stdout, " RAFT ", log.LstdFlags)
	config.LocalID = raft.ServerID(args.name)

	return config
}

func logStore(args *Arguments) raft.LogStore {
	if args.persist() {
		store, _ := raftboltdb.NewBoltStore(fmt.Sprintf("%s/logs.dat", args.dataDirectory))
		return store
	}
	return raft.NewInmemStore()
}

func stableStore(args *Arguments) raft.StableStore {
	if args.persist() {
		store, _ := raftboltdb.NewBoltStore(fmt.Sprintf("%s/stable.dat", args.dataDirectory))
		return store
	}
	return raft.NewInmemStore()
}

func snapStore(args *Arguments) raft.SnapshotStore {
	if args.persist() {
		store, _ := raft.NewFileSnapshotStoreWithLogger(args.dataDirectory, 5, logger)
		return store
	}
	return raft.NewInmemSnapshotStore()
}

func leadershipChanges(channel <-chan bool) {
	for m := range channel {
		isCurrentLeader = m
	}
}

func evicted(key, value interface{}) {
	logger.Printf("EVICTED: %+v/%+v", key, value)
}

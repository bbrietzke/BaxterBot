package store

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

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
	config.Logger = log.New(os.Stdout, "RAFT  ", log.LstdFlags)
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

func bootstrapOrJoin(args *Arguments, myAddr raft.ServerAddress) {
	time.Sleep(defaultTimeOut)
	if c := repl.GetConfiguration(); c.Error() == nil {
		if len(c.Configuration().Servers) == 0 && len(args.join) == 0 {
			logger.Println("* * * * * * * * * * BOOTSTRAPPING * * * * * * * * * * ")
			if err := repl.BootstrapCluster(raft.Configuration{Servers: []raft.Server{{ID: raft.ServerID(args.name), Address: myAddr}}}).Error(); err != nil {
				logger.Fatalln(err)
			}
		} else if len(args.join) > 0 {
			logger.Println("Should connect to a different server here")
		}
	}
}

func evicted(key, value interface{}) {
	logger.Printf("EVICTED: %+v/%+v", key, value)
}

package store

import (
	"time"

	"github.com/hashicorp/raft"
)

func bootstrapOrJoin(args *Arguments, myAddr raft.ServerAddress) {
	time.Sleep(defaultTimeOut)
	config := repl.GetConfiguration()

	if config.Error() != nil {
		logger.Println(config.Error())
	}

	cnt := len(config.Configuration().Servers)

	if cnt == 0 && len(args.join) == 0 {
		logger.Println("* * * * * * * * * * BOOTSTRAPPING * * * * * * * * * * ")
		if err := repl.BootstrapCluster(raft.Configuration{Servers: []raft.Server{{ID: raft.ServerID(args.name), Address: myAddr}}}).Error(); err != nil {
			logger.Fatalln(err)
		}
	}

	if cnt == 0 && len(args.join) > 0 {
		logger.Println("* * * * * * * * * * * SEARCHING * * * * * * * * * * * ")
	}
}

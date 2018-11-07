package swarm

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/bbrietzke/BaxterBot/pkg/protocol"
)

func isLeader() bool {
	future := swarmer.VerifyLeader()
	err := future.Error()
	return err == nil
}

func leaderAddr() string {
	return string(swarmer.Leader())
}

func apply(value []byte) error {
	return swarmer.Apply(value, applyTimeout).Error()
}

func monitorLeadership(ctx context.Context) {
	t := time.NewTicker(monitorLeadershipTimeout)

	for {
		select {
		case <-ctx.Done():
			logger.Println("MonitorLeadership stopped")
			return
		case <-t.C:
			if isLeader() {
				updateMyHTTP()
			}
		}
	}
}

func listenForLeadership(c <-chan bool) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for msg := range c {
		if msg && !leader {
			logger.Println("* * * * * * * * * * LEADERSHIP ASSUMED * * * * * * * * * *")
			updateMyHTTP()
			ctx, cancel = context.WithCancel(context.Background())
			defer cancel()
			go monitorLeadership(ctx)
		} else if !msg && leader {
			cancel()
		}

		leader = msg
	}
}

func outboundIP(port string) string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String() + port
}

func updateMyHTTP() {
	value := outboundIP(httpPort)
	v, err := protocol.UpdateLeaderHTTP(value)
	if err != nil {
		logger.Print(err)
	}

	err = apply(v)
	if err != nil {
		logger.Print(err)
	}
}

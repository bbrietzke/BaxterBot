package swarm

import (
	"context"
	"log"
	"net"
	"time"
)

func isLeader() bool {
	future := swarmer.VerifyLeader()
	err := future.Error()
	return err == nil
}

func leaderAddr() string {
	return string(swarmer.Leader())
}

func monitorLeadership(ctx context.Context) {
	t := time.NewTicker(600 * time.Second)

	for {
		select {
		case <-ctx.Done():
			logger.Println("MonitorLeadership stopped")
			return
		case <-t.C:
			if isLeader() {
				if err := swarmer.Apply(updateLeaderHTTP(), 3*time.Second).Error(); err != nil {
					logger.Println("MonitorLeadership failed")
				}
			}
		}
	}
}

func listenForLeadership(c <-chan bool) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for msg := range c {
		logger.Println(msg, leader)
		if msg && !leader {
			logger.Println("* * * * * * * * * * LEADERSHIP ASSUMED * * * * * * * * * *")
			swarmer.Apply(updateLeaderHTTP(), 3*time.Second).Error()
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

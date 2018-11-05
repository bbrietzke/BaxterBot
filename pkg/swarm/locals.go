package swarm

import (
	"log"
	"net"
)

func isLeader() bool {
	return leader
}

func leaderAddr() string {
	return string(swarmer.Leader())
}

func listenForLeadership(c <-chan bool) {
	for msg := range c {
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

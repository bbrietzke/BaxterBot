package swarm

func isLeader() bool {
	return leader
}

func listenForLeadership(c <-chan bool) {
	for msg := range c {
		leader = msg
	}
}

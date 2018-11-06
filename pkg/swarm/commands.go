package swarm

import (
	"encoding/json"
)

type commandIndex int16

const (
	leaderUpdate commandIndex = iota
)

type command struct {
	Idx commandIndex    `json:"idx"`
	Sub json.RawMessage `json:"sub,omitempty"`
}

type leaderHTTP struct {
	Addr string `json:"addr"`
}

func updateLeaderHTTP() []byte {
	a, _ := json.Marshal(leaderHTTP{Addr: outboundIP(httpPort)})
	c := command{
		Idx: leaderUpdate,
		Sub: a,
	}
	d, _ := json.Marshal(c)

	return d
}

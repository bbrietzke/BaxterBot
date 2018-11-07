package swarm

import (
	"encoding/json"
)

type commandType int16

const (
	cmdUpdateLeaderAddr commandType = iota
	cmdUpdateKeyValuePair
)

type command struct {
	Typ commandType     `json:"typ"`
	Sub json.RawMessage `json:"sub,omitempty"`
}

type keyValueUpdateJSON struct {
	Key   interface{} `json:"key,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

type leaderHTTPJSON struct {
	Addr string `json:"addr"`
}

func updateKeyValue(key, value interface{}) []byte {
	a, _ := json.Marshal(keyValueUpdateJSON{Key: key, Value: value})
	c := command{
		Idx: keyValueUpdate,
		Sub: a,
	}
	d, _ := json.Marshal(c)

	return d
}

func updateLeaderHTTP() []byte {
	a, _ := json.Marshal(leaderHTTPJSON{Addr: outboundIP(httpPort)})
	c := command{
		Idx: leaderUpdate,
		Sub: a,
	}
	d, _ := json.Marshal(c)

	return d
}

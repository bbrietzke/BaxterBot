package swarm

import (
	"encoding/json"
	"errors"
	"io"
	"sync"

	"github.com/hashicorp/raft"
)

func newStateMachine() raft.FSM {
	return &stateMachine{}
}

type stateMachine struct {
	sync.Mutex
	masters map[string]struct{}
}

func (fsm *stateMachine) Apply(l *raft.Log) interface{} {
	cmd := command{}
	if err := json.Unmarshal(l.Data, &cmd); err == nil {
		switch cmd.Idx {
		case keyValueUpdate:
			d := keyValueUpdateJSON{}
			json.Unmarshal(cmd.Sub, &d)
			logger.Printf("index %d :: data %+v", l.Index, d)
		case leaderUpdate:
			d := leaderHTTPJSON{}
			json.Unmarshal(cmd.Sub, &d)
			logger.Printf("index %d :: data %+v", l.Index, d)
			leaderNetAddr = d.Addr
		}
	}
	return nil
}

func (fsm *stateMachine) Restore(snap io.ReadCloser) error {
	return errors.New("Not Implemented")
}

func (fsm *stateMachine) Snapshot() (raft.FSMSnapshot, error) {
	return nil, errors.New("Not Implemented")
}

func (fsm *stateMachine) Persist(sink raft.SnapshotSink) error {
	return nil
}

func (fsm *stateMachine) Release() {

}

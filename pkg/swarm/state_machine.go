package swarm

import (
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
	logger.Panicf("%+v", l)
	return nil
}

func (fsm *stateMachine) Restore(snap io.ReadCloser) error {
	return errors.New("Not Implemented")
}

func (fsm *stateMachine) Snapshot() (raft.FSMSnapshot, error) {
	return nil, errors.New("Not Implemented")
}

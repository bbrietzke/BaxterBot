package store

import (
	"io"

	"github.com/hashicorp/raft"
	"github.com/pkg/errors"
)

func newStateMachine(args *Arguments) raft.FSM {
	return &stateMachine{}
}

type stateMachine struct {
}

func (fsm *stateMachine) Apply(l *raft.Log) interface{} {
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

func (fsm *stateMachine) Release() {}

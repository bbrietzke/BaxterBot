package store

import (
	"io"

	"github.com/bbrietzke/BaxterBot/pkg/protocol"
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/raft"
	"github.com/pkg/errors"
)

func newStateMachine(args *Arguments) raft.FSM {
	return &stateMachine{}
}

type stateMachine struct {
}

func (fsm *stateMachine) Apply(l *raft.Log) interface{} {
	wrapper := &protocol.CommandWrapper{}
	proto.Unmarshal(l.Data, wrapper)

	switch wrapper.Type {
	case protocol.CommandWrapper_KEY_VALUE_DELETE:
	case protocol.CommandWrapper_KEY_VALUE_CREATE:
	default:
		logger.Println(wrapper.Type)
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

func (fsm *stateMachine) Release() {}

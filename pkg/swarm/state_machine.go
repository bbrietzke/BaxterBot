package swarm

import (
	"errors"
	"io"
	"sync"

	"github.com/golang/protobuf/proto"

	"github.com/bbrietzke/BaxterBot/pkg/protocol"

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
	wrapper := &protocol.CommandWrapper{}
	proto.Unmarshal(l.Data, wrapper)

	switch wrapper.Type {
	case protocol.CommandWrapper_KEY_VALUE_DELETE:
		fsm.deleteKeyValue(wrapper.GetCmd())
	case protocol.CommandWrapper_KEY_VALUE_CREATE:
		fsm.createKeyValue(wrapper.GetCmd())
	case protocol.CommandWrapper_LEADER_HTTP_UPDATE:
		fsm.updateLeaderHTTP(wrapper.GetCmd())
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

func (fsm *stateMachine) updateLeaderHTTP(v []byte) {
	value := protocol.LeaderHttpUpdate{}
	err := proto.Unmarshal(v, &value)
	if err != nil {
		logger.Println(err)
	}
	leaderNetAddr = value.Addr
	logger.Println("leader http has been set to", leaderNetAddr)
	outputChannel <- value
}

func (fsm *stateMachine) createKeyValue(v []byte) {
	value := protocol.KeyValueCreate{}
	err := proto.Unmarshal(v, &value)
	if err != nil {
		logger.Println(err)
	}
	outputChannel <- value
}

func (fsm *stateMachine) deleteKeyValue(v []byte) {
	value := protocol.KeyValueDelete{}
	err := proto.Unmarshal(v, &value)
	if err != nil {
		logger.Println(err)
	}
	outputChannel <- value
}

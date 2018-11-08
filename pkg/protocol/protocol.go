package protocol

import (
	proto "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
)

// CommandPipeline defines a channel type for sending and receiving commands
type CommandPipeline chan interface{}

// CommandPipelineReceiver accepts commands that were put into the Pipeline
type CommandPipelineReceiver <-chan interface{}

const (
	pipelineBuffer int32 = 5
)

var (
	cmdPipe chan interface{}
)

func init() {
	cmdPipe = make(chan interface{}, pipelineBuffer)
}

// Pipeline returns
func Pipeline() CommandPipelineReceiver {
	return cmdPipe
}

// UpdateLeaderHTTP creates the byte stream to update the leader http ip and port.
func UpdateLeaderHTTP(leaderAddr string) ([]byte, error) {
	b, e := proto.Marshal(&LeaderHttpUpdate{Addr: leaderAddr})
	return wrapCommand(CommandWrapper_LEADER_HTTP_UPDATE, b, e)
}

// CreateKeyValuePair creates the byte stream for adding key value pairs.
func CreateKeyValuePair(key string, value []byte) ([]byte, error) {
	b, e := proto.Marshal(&KeyValueCreate{Key: key, Value: &any.Any{Value: value}})
	return wrapCommand(CommandWrapper_KEY_VALUE_CREATE, b, e)
}

// DeleteKeyValuePair creates the byte stream for removing key/value pairs.
func DeleteKeyValuePair(key string) ([]byte, error) {
	b, e := proto.Marshal(&KeyValueDelete{Key: key})
	return wrapCommand(CommandWrapper_KEY_VALUE_DELETE, b, e)
}

func wrapCommand(v CommandWrapper_CommandType, value []byte, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}

	return proto.Marshal(&CommandWrapper{Type: v, Child: &any.Any{Value: value}})
}

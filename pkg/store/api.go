package store

import (
	"bytes"
	"encoding/gob"

	"github.com/bbrietzke/BaxterBot/pkg/protocol"

	"github.com/pkg/errors"
)

// Add will create a new entry in the store with the given value.
// This method will throw an error if it is not performed on the replication master.
func Add(key string, value interface{}) error {
	if !isCurrentLeader {
		return errors.New("not the replication master")
	}

	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(value)
	if err != nil {
		return errors.Wrap(err, "encoding value")
	}
	v, err := protocol.CreateKeyValuePair(key, encoded.Bytes())
	if err != nil {
		return errors.Wrap(err, "CreateKeyValuePair")
	}
	if f := repl.Apply(v, defaultTimeOut).Error(); f != nil {
		return errors.Wrap(f, "Apply")
	}
	return nil
}

// Delete will create a new entry in the store with the given value.
// This method will throw an error if it is not performed on the replication master.
func Delete(key string) error {
	if !isCurrentLeader {
		return errors.New("not the replication master")
	}

	v, err := protocol.DeleteKeyValuePair(key)
	if err != nil {
		return errors.Wrap(err, "DeleteKeyValuePair")
	}
	if f := repl.Apply(v, defaultTimeOut).Error(); f != nil {
		return errors.Wrap(f, "Apply")
	}
	return nil
}

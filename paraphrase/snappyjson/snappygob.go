// Package snappyjson implements snappy compression of JSON objects for StormDB
package snappyjson

import (
	msgpack "gopkg.in/vmihailenco/msgpack.v2"

	"github.com/golang/snappy"
)

const msgpackName = "snappymsgpack"

// Codec that encodes to and decodes from JSON then is compressed with snappy.
var MsgpackCodec = new(snappyMsgpackCodec)

type snappyMsgpackCodec int

func (j snappyMsgpackCodec) Marshal(v interface{}) ([]byte, error) {
	bytes, err := msgpack.Marshal(v)

	if err != nil {
		return nil, err
	}

	return snappy.Encode(nil, bytes), nil
}

func (j snappyMsgpackCodec) Unmarshal(b []byte, v interface{}) error {
	bytes, err := snappy.Decode(nil, b)

	if err != nil {
		return err
	}

	return msgpack.Unmarshal(bytes, v)
}

func (j snappyMsgpackCodec) Name() string {
	return msgpackName
}

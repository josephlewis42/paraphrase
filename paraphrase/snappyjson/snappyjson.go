// Package snappyjson implements snappy compression of JSON objects for StormDB
package snappyjson

import (
	"encoding/json"

	"github.com/golang/snappy"
)

const name = "snappyjson"

// Codec that encodes to and decodes from JSON then is compressed with snappy.
var Codec = new(snappyJsonCodec)

type snappyJsonCodec int

func (j snappyJsonCodec) Marshal(v interface{}) ([]byte, error) {
	bytes, err := json.Marshal(v)

	if err != nil {
		return nil, err
	}

	return snappy.Encode(nil, bytes), nil
}

func (j snappyJsonCodec) Unmarshal(b []byte, v interface{}) error {
	bytes, err := snappy.Decode(nil, b)

	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, v)
}

func (j snappyJsonCodec) Name() string {
	return name
}

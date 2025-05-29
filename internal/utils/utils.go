package utils

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

func UnmarshalJSON[T any](data []byte) (T, error) {
	var target T
	err := json.Unmarshal(data, &target)
	return target, err
}

func UnmarshalGob[T any](data []byte) (T, error) {
	buff := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buff)
	var target T
	err := decoder.Decode(&target)
	if err != nil {
		return target, err
	}
	return target, err
}

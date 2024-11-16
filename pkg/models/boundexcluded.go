package models

import (
	"reflect"

	"github.com/fxamacker/cbor/v2"

	"github.com/tai-kun/surrealdb.go/pkg/codec"
)

type BoundExcluded[T any] struct {
	Value T
}

func NewBoundExcluded[T any](v T) *BoundExcluded[T] {
	return &BoundExcluded[T]{
		Value: v,
	}
}

func (be *BoundExcluded[T]) MarshalCBOR() ([]byte, error) {
	return CBORFormatter.Marshal(cbor.Tag{
		Number:  TagBoundExcluded,
		Content: be.Value,
	})
}

func (be *BoundExcluded[T]) UnmarshalCBOR(data []byte) error {
	return be.unmarshal(CBORFormatter, data)
}

func (be *BoundExcluded[T]) MarshalJSON() ([]byte, error) {
	return JSONFormatter.Marshal(be.Value)
}

func (be *BoundExcluded[T]) UnmarshalJSON(data []byte) error {
	return be.unmarshal(JSONFormatter, data)
}

func (be *BoundExcluded[T]) unmarshal(f codec.Formatter, data []byte) error {
	var t T
	if err := f.Unmarshal(data, &t); err != nil {
		return err
	}

	be.Value = t
	return nil
}

func (be *BoundExcluded[T]) SurrealString() (string, error) {
	if reflect.TypeOf(be.Value) == reflect.TypeOf(None{}) {
		return "", nil
	}

	j, err := be.MarshalJSON()
	if err != nil {
		return "", err
	}

	s := string(j)
	if s == "null" {
		s = ""
	}

	return s, nil
}

func (be *BoundExcluded[T]) value() T {
	return be.Value
}

func (be *BoundExcluded[T]) excluded() bool {
	return true
}

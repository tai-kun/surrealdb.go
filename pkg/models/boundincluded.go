package models

import (
	"github.com/fxamacker/cbor/v2"

	"github.com/tai-kun/surrealdb.go/pkg/codec"
)

type BoundIncluded[T any] struct {
	Value T
}

func NewBoundIncluded[T any](v T) *BoundIncluded[T] {
	return &BoundIncluded[T]{
		Value: v,
	}
}

func (bi *BoundIncluded[T]) MarshalCBOR() ([]byte, error) {
	return CBORFormatter.Marshal(cbor.Tag{
		Number:  TagBoundIncluded,
		Content: bi.Value,
	})
}

func (bi *BoundIncluded[T]) UnmarshalCBOR(data []byte) error {
	return bi.unmarshal(CBORFormatter, data)
}

func (bi *BoundIncluded[T]) MarshalJSON() ([]byte, error) {
	return JSONFormatter.Marshal(bi.Value)
}

func (bi *BoundIncluded[T]) UnmarshalJSON(data []byte) error {
	return bi.unmarshal(JSONFormatter, data)
}

func (bi *BoundIncluded[T]) unmarshal(f codec.Formatter, data []byte) error {
	var t T
	if err := f.Unmarshal(data, &t); err != nil {
		return err
	}

	bi.Value = t
	return nil
}

func (bi *BoundIncluded[T]) SurrealString() (string, error) {
	j, err := bi.MarshalJSON()
	if err != nil {
		return "", err
	}

	s := string(j)
	if s == "null" {
		s = ""
	}

	return s, nil
}

func (bi *BoundIncluded[T]) value() T {
	return bi.Value
}

func (bi *BoundIncluded[T]) excluded() bool {
	return false
}

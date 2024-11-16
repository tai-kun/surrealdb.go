package models

import "fmt"

type Bound[T any] struct {
	ex    bool
	Value T
}

func (b *Bound[T]) Excluded() bool {
	return b.ex
}

func (b *Bound[T]) Included() bool {
	return !b.ex
}

func (b *Bound[T]) Exclude() {
	b.ex = true
}

func (b *Bound[T]) Include() {
	b.ex = false
}

func (b *Bound[T]) MarshalCBOR() ([]byte, error) {
	return b.toSpecificBound().MarshalCBOR()
}

func (b *Bound[T]) UnmarshalCBOR(data []byte) error {
	err := fmt.Errorf("cannot unmarshal directly using CBOR formatter")
	return err
}

func (b *Bound[T]) MarshalJSON() ([]byte, error) {
	return b.toSpecificBound().MarshalJSON()
}

func (b *Bound[T]) UnmarshalJSON(data []byte) error {
	err := fmt.Errorf("cannot unmarshal directly using JSON formatter")
	return err
}

func (b *Bound[T]) SurrealString() (string, error) {
	return b.toSpecificBound().SurrealString()
}

func (b *Bound[T]) toSpecificBound() interface {
	MarshalCBOR() ([]byte, error)
	MarshalJSON() ([]byte, error)
	SurrealString() (string, error)
} {
	if b.Excluded() {
		return NewBoundExcluded(b.Value)
	} else {
		return NewBoundIncluded(b.Value)
	}
}

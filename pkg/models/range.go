package models

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

type Range[T any] struct {
	Begin *Bound[T]
	End   *Bound[T]
}

type bound[T any] interface {
	value() T
	excluded() bool
}

func NewRange[T any](begin bound[T], end bound[T]) *Range[T] {
	var (
		b *Bound[T] = nil
		e *Bound[T] = nil
	)
	if begin != nil {
		b = &Bound[T]{
			ex:    begin.excluded(),
			Value: begin.value(),
		}
	}
	if end != nil {
		e = &Bound[T]{
			ex:    end.excluded(),
			Value: end.value(),
		}
	}
	return &Range[T]{
		Begin: b,
		End:   e,
	}
}

func (r *Range[T]) MarshalCBOR() ([]byte, error) {
	return CBORFormatter.Marshal(cbor.Tag{
		Number:  TagRange,
		Content: [2]any{r.Begin, r.End},
	})
}

func (r *Range[T]) UnmarshalCBOR(data []byte) (err error) {
	var c [2]*cbor.RawTag
	if err := CBORFormatter.Unmarshal(data, &c); err != nil {
		return err
	}

	r.Begin, err = reviveBound[T](c[0])
	if err != nil {
		return err
	}

	r.End, err = reviveBound[T](c[1])
	if err != nil {
		return err
	}

	return nil
}

func reviveBound[T any](rt *cbor.RawTag) (*Bound[T], error) {
	if rt == nil {
		return nil, nil
	}

	switch rt.Number {
	case TagBoundExcluded:
		var be BoundExcluded[T]
		if err := be.UnmarshalCBOR(rt.Content); err != nil {
			return nil, err
		}
		return &Bound[T]{ex: true, Value: be.Value}, nil

	case TagBoundIncluded:
		var bi BoundIncluded[T]
		if err := bi.UnmarshalCBOR(rt.Content); err != nil {
			return nil, err
		}
		return &Bound[T]{ex: false, Value: bi.Value}, nil

	default:
		return nil, fmt.Errorf("invalid tag=%d", rt.Number)
	}
}

func (r *Range[T]) MarshalJSON() ([]byte, error) {
	s, err := r.SurrealString()
	if err != nil {
		return nil, err
	}
	return JSONFormatter.Marshal(s)
}

func (r *Range[T]) UnmarshalJSON(data []byte) error {
	err := fmt.Errorf("cannot unmarshal using JSON formatter")
	return err
}

func (r *Range[T]) SurrealString() (string, error) {
	s := ""

	if r.Begin != nil {
		if r.Begin.Excluded() {
			s += ">"
		}

		b, err := r.Begin.SurrealString()
		if err != nil {
			return "", err
		}

		s = b + s
	}

	s += ".."

	if r.End != nil {
		if r.End.Included() {
			s += "="
		}

		e, err := r.End.SurrealString()
		if err != nil {
			return "", err
		}

		s += e
	}

	return s, nil
}

package models

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/tai-kun/surrealdb.go/pkg/utils"
)

type RecordID[T any] struct {
	Table string
	ID    T
}

func NewRecordID[T any](table string, id T) *RecordID[T] {
	return &RecordID[T]{
		Table: table,
		ID:    id,
	}
}

func (r *RecordID[T]) MarshalCBOR() ([]byte, error) {
	return CBORFormatter.Marshal(cbor.Tag{
		Number:  TagRecordID,
		Content: [2]any{r.Table, r.ID},
	})
}

func (r *RecordID[T]) UnmarshalCBOR(data []byte) error {
	var c [2]cbor.RawMessage
	if err := CBORFormatter.Unmarshal(data, &c); err != nil {
		return err
	}

	var t string
	if err := CBORFormatter.Unmarshal(c[0], &t); err != nil {
		return err
	}

	var i T
	if err := CBORFormatter.Unmarshal(c[1], &i); err != nil {
		return err
	}

	r.Table = t
	r.ID = i
	return nil
}

func (r *RecordID[T]) MarshalJSON() ([]byte, error) {
	s, err := r.string()
	if err != nil {
		return nil, err
	}

	return JSONFormatter.Marshal(s)
}

func (r *RecordID[T]) UnmarshalJSON(data []byte) error {
	err := fmt.Errorf("cannot unmarshal using JSON formatter")
	return err
}

func (r *RecordID[T]) SurrealString() (string, error) {
	s, err := r.string()
	if err != nil {
		return "", err
	}

	return "r" + utils.QuoteStr(s), nil
}

func (r *RecordID[T]) string() (string, error) {
	type SurrealStringer interface {
		SurrealString() (string, error)
	}
	if v, ok := any(r.ID).(SurrealStringer); ok {
		i, err := v.SurrealString()
		if err != nil {
			return "", err
		}

		return utils.QuoteRID(r.Table) + ":" + i, nil
	}

	i, err := JSONFormatter.Marshal(r.ID)
	if err != nil {
		return "", err
	}

	if len(i) > 0 && i[0] == '-' {
		return utils.QuoteRID(r.Table) + ":" + utils.QuoteRID(string(i)), nil
	}
	return utils.QuoteRID(r.Table) + ":" + string(i), nil
}

package models

import (
	"github.com/fxamacker/cbor/v2"
	"github.com/tai-kun/surrealdb.go/pkg/utils"
)

type Table string

func NewTable(v string) Table {
	return Table(v)
}

func (t Table) MarshalCBOR() ([]byte, error) {
	return CBORFormatter.Marshal(cbor.Tag{
		Number:  TagTable,
		Content: string(t),
	})
}

func (t *Table) UnmarshalCBOR(data []byte) error {
	var c string
	if err := CBORFormatter.Unmarshal(data, &c); err != nil {
		return err
	}

	*t = Table(c)
	return nil
}

func (t Table) MarshalJSON() ([]byte, error) {
	return JSONFormatter.Marshal(string(t))
}

func (t *Table) UnmarshalJSON(data []byte) error {
	var c string
	if err := JSONFormatter.Unmarshal(data, &c); err != nil {
		return err
	}

	*t = Table(c)
	return nil
}

func (t Table) SurrealString() (string, error) {
	// SurrealDB では escape_ident でエスケープしている:
	// https://github.com/surrealdb/surrealdb/blob/v2.0.4/core/src/sql/table.rs#L78
	return utils.QuoteIdent(string(t)), nil
}

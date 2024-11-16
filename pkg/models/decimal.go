package models

import (
	"github.com/fxamacker/cbor/v2"
)

type Decimal string

func NewDecimal(v string) Decimal {
	return Decimal(v)
}

func (d Decimal) MarshalCBOR() ([]byte, error) {
	return CBORFormatter.Marshal(cbor.Tag{
		Number:  TagDecimal,
		Content: string(d),
	})
}

func (d *Decimal) UnmarshalCBOR(data []byte) error {
	var c string
	if err := CBORFormatter.Unmarshal(data, &c); err != nil {
		return err
	}

	*d = Decimal(c)
	return nil
}

func (d Decimal) MarshalJSON() ([]byte, error) {
	return JSONFormatter.Marshal(string(d))
}

func (d *Decimal) UnmarshalJSON(data []byte) error {
	var c string
	if err := JSONFormatter.Unmarshal(data, &c); err != nil {
		return err
	}

	*d = Decimal(c)
	return nil
}

func (d Decimal) SurrealString() (string, error) {
	return string(d) + "dec", nil
}

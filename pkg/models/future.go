package models

import (
	"github.com/fxamacker/cbor/v2"
)

type Future string

func NewFuture(v string) Future {
	return Future(v)
}

func (f Future) MarshalCBOR() ([]byte, error) {
	return CBORFormatter.Marshal(cbor.Tag{
		Number:  TagFuture,
		Content: string(f),
	})
}

func (f *Future) UnmarshalCBOR(data []byte) error {
	var c string
	if err := CBORFormatter.Unmarshal(data, &c); err != nil {
		return err
	}

	*f = Future(c)
	return nil
}

func (f Future) MarshalJSON() ([]byte, error) {
	return JSONFormatter.Marshal(string(f))
}

func (f *Future) UnmarshalJSON(data []byte) error {
	var c string
	if err := JSONFormatter.Unmarshal(data, &c); err != nil {
		return err
	}

	*f = Future(c)
	return nil
}

func (f Future) SurrealString() (string, error) {
	return "<future>{" + string(f) + "}", nil
}

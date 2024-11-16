package models

import (
	"time"

	"github.com/fxamacker/cbor/v2"
)

const DatetimeLayout = "2006-01-02T15:04:05.000000000Z"

type Datetime struct {
	time.Time
}

func NewDatetime() *Datetime {
	return &Datetime{time.Now()}
}

func (d *Datetime) MarshalCBOR() ([]byte, error) {
	nsTime := d.UnixNano()
	s := nsTime / 1_000_000_000
	ns := nsTime % 1_000_000_000

	return CBORFormatter.Marshal(cbor.Tag{
		Number:  TagDatetime,
		Content: [2]int64{s, ns},
	})
}

func (d *Datetime) UnmarshalCBOR(data []byte) error {
	var c [2]int64
	if err := CBORFormatter.Unmarshal(data, &c); err != nil {
		return err
	}

	*d = Datetime{time.Unix(c[0], c[1])}
	return nil
}

func (d Datetime) MarshalJSON() ([]byte, error) {
	return JSONFormatter.Marshal(d.UTC().Format(DatetimeLayout))
}

func (d *Datetime) SurrealString() (string, error) {
	return "d'" + d.UTC().Format(DatetimeLayout) + "'", nil
}

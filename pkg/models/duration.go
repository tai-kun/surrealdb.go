package models

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/fxamacker/cbor/v2"
)

type Duration time.Duration

func NewDuration(v time.Duration) Duration {
	return Duration(v)
}

func (d Duration) MarshalCBOR() ([]byte, error) {
	nsTime := int64(d)
	s := nsTime / 1_000_000_000
	ns := nsTime % 1_000_000_000

	return CBORFormatter.Marshal(cbor.Tag{
		Number:  TagDuration,
		Content: [2]int64{s, ns},
	})
}

func (d *Duration) UnmarshalCBOR(data []byte) error {
	var c [2]int64
	if err := CBORFormatter.Unmarshal(data, &c); err != nil {
		return err
	}

	*d = Duration(c[0]*1_000_000_000 + c[1])
	return nil
}

func (d Duration) MarshalJSON() ([]byte, error) {
	s, err := d.SurrealString()
	if err != nil {
		return nil, err
	}
	return JSONFormatter.Marshal(s)
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	var c string
	if err := JSONFormatter.Unmarshal(data, &c); err != nil {
		return err
	}

	pd, err := ParseDuration(c)
	if err != nil {
		return err
	}

	*d = pd
	return nil
}

func (d Duration) SurrealString() (string, error) {
	if d <= 0 {
		return "0ns", nil
	}

	const (
		SECONDS_PER_MINUTE          uint64 = 60
		SECONDS_PER_HOUR            uint64 = 60 * SECONDS_PER_MINUTE
		SECONDS_PER_DAY             uint64 = 24 * SECONDS_PER_HOUR
		SECONDS_PER_WEEK            uint64 = 7 * SECONDS_PER_DAY
		SECONDS_PER_YEAR            uint64 = 365 * SECONDS_PER_DAY
		NANOSECONDS_PER_MICROSECOND uint64 = 1_000
		NANOSECONDS_PER_MILLISECOND uint64 = 1_000_000
		NANOSECONDS_PER_SECOND      uint64 = 1_000_000_000
	)
	var (
		s    = ""
		t    = uint64(d)
		year uint64
		week uint64
		days uint64
		hour uint64
		mins uint64
		secs uint64 = t / NANOSECONDS_PER_SECOND
		msec uint64
		usec uint64
		nano uint64 = t % NANOSECONDS_PER_SECOND
	)
	year = secs / SECONDS_PER_YEAR
	secs = secs % SECONDS_PER_YEAR
	week = secs / SECONDS_PER_WEEK
	secs = secs % SECONDS_PER_WEEK
	days = secs / SECONDS_PER_DAY
	secs = secs % SECONDS_PER_DAY
	hour = secs / SECONDS_PER_HOUR
	secs = secs % SECONDS_PER_HOUR
	mins = secs / SECONDS_PER_MINUTE
	secs = secs % SECONDS_PER_MINUTE
	msec = nano / NANOSECONDS_PER_MILLISECOND
	nano = nano % NANOSECONDS_PER_MILLISECOND
	usec = nano / NANOSECONDS_PER_MICROSECOND
	nano = nano % NANOSECONDS_PER_MICROSECOND
	if year > 0 {
		s += strconv.FormatUint(year, 10) + "y"
	}
	if week > 0 {
		s += strconv.FormatUint(week, 10) + "w"
	}
	if days > 0 {
		s += strconv.FormatUint(days, 10) + "d"
	}
	if hour > 0 {
		s += strconv.FormatUint(hour, 10) + "h"
	}
	if mins > 0 {
		s += strconv.FormatUint(mins, 10) + "m"
	}
	if secs > 0 {
		s += strconv.FormatUint(secs, 10) + "s"
	}
	if msec > 0 {
		s += strconv.FormatUint(msec, 10) + "ms"
	}
	if usec > 0 {
		s += strconv.FormatUint(usec, 10) + "µs"
	}
	if nano > 0 {
		s += strconv.FormatUint(nano, 10) + "ns"
	}

	return s, nil
}

func ParseDuration(s string) (Duration, error) {
	if s == "0" || s == "0ns" {
		return 0, nil
	}
	if s == "" {
		err := errors.New("surrealdb: models: invalid duration \"\": empty")
		return 0, err
	}

	const (
		SECONDS_PER_MINUTE          uint64 = 60
		SECONDS_PER_HOUR            uint64 = 60 * SECONDS_PER_MINUTE
		SECONDS_PER_DAY             uint64 = 24 * SECONDS_PER_HOUR
		SECONDS_PER_WEEK            uint64 = 7 * SECONDS_PER_DAY
		SECONDS_PER_YEAR            uint64 = 365 * SECONDS_PER_DAY
		NANOSECONDS_PER_MICROSECOND uint64 = 1_000
		NANOSECONDS_PER_MILLISECOND uint64 = 1_000_000
		NANOSECONDS_PER_SECOND      uint64 = 1_000_000_000
	)
	var (
		orig = s
		secs uint64
		nano uint64
	)
	for s != "" {
		i := 0
		for ; i < len(s) && s[i] >= '0' && s[i] <= '9'; i++ {
		}
		if i == 0 {
			err := errors.New(
				"surrealdb: models: invalid duration " + strconv.Quote(orig) + ": no value",
			)
			return 0, err
		}

		v, err := strconv.ParseUint(s[:i], 10, 64)
		if err != nil {
			err := fmt.Errorf(
				"surrealdb: models: invalid duration %s: %w",
				strconv.Quote(orig), err,
			)
			return 0, err
		}

		s = s[i:]
		i = 0
		for ; i < len(s) && (s[i] < '0' || s[i] > '9'); i++ {
		}

		switch s[:i] {
		case "y":
			secs += v * SECONDS_PER_YEAR
		case "w":
			secs += v * SECONDS_PER_WEEK
		case "d":
			secs += v * SECONDS_PER_DAY
		case "h":
			secs += v * SECONDS_PER_HOUR
		case "m":
			secs += v * SECONDS_PER_MINUTE
		case "s":
			secs += v
		case "ms":
			nano += v * NANOSECONDS_PER_MILLISECOND
		case "us", "µs", "μs":
			nano += v * NANOSECONDS_PER_MICROSECOND
		case "ns":
			nano += v
		default:
			err := errors.New(
				"surrealdb: models: invalid duration " + strconv.Quote(orig) +
					": invaid unit " + strconv.Quote(s[:i]),
			)
			return 0, err
		}
		s = s[i:]
	}

	return Duration(secs*1_000_000_000 + nano), nil
}

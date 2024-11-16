package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tai-kun/surrealdb.go/pkg/models"
)

func TestDurationSurrealString(t *testing.T) {
	tests := map[string]models.Duration{
		"0ns":       0,
		"12µs345ns": 12_345,
	}
	for expected, d := range tests {
		if actual, err := d.SurrealString(); assert.NoError(t, err) {
			assert.Equal(t, expected, actual)
		}
	}
}

func TestDurationCBOR(t *testing.T) {
	tests := map[string]models.Duration{
		"0ns":       0,
		"12µs345ns": 12_345,
	}
	for expected, src := range tests {
		data, err := models.CBORFormatter.Marshal(src)
		if assert.NoError(t, err) {
			var dst models.Duration
			if err := models.CBORFormatter.Unmarshal(data, &dst); assert.NoError(t, err) {
				if actual, err := dst.SurrealString(); assert.NoError(t, err) {
					assert.Equal(t, expected, actual)
				}
			}
		}
	}
}

func TestDurationJSON(t *testing.T) {
	tests := map[string]models.Duration{
		"0ns":       0,
		"12µs345ns": 12_345,
	}
	for expected, src := range tests {
		data, err := models.JSONFormatter.Marshal(src)
		if assert.NoError(t, err) {
			var dst string
			if err := models.JSONFormatter.Unmarshal(data, &dst); assert.NoError(t, err) {
				assert.Equal(t, expected, dst)
			}
		}
	}
}

func TestParseDuration(t *testing.T) {
	tests := map[string]string{
		"0":             "0ns",
		"12345ns":       "12µs345ns",
		"012345us":      "12ms345µs",
		"012345µs":      "12ms345µs",
		"012345μs":      "12ms345µs",
		"123456ms":      "2m3s456ms",
		"3601s":         "1h1s",
		"10000m":        "6d22h40m",
		"90d":           "12w6d",
		"100w":          "1y47w6d",
		"2y":            "2y",
		"1ns1us1ms1s1m": "1m1s1ms1µs1ns",
		"1ms1m":         "1m1ms",
		"1m1s1m1s1m1s":  "3m3s",
	}
	for s, expected := range tests {
		if d, err := models.ParseDuration(s); assert.NoError(t, err) {
			if actual, err := d.SurrealString(); err != nil {
				assert.Equal(t, expected, actual)
			}
		}
	}
}

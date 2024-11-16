package models_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/tai-kun/surrealdb.go/pkg/models"
)

func TestDatetimeSurrealString(t *testing.T) {
	tests := map[string]models.Datetime{
		"d'2024-06-01T12:34:56.780123456Z'": {time.Unix(1717245296, 780123456)},
		// "d'275760-09-13T00:00:00.000000000Z'":  {time.Unix(8640000000000, 0)},
		// "d'-271821-04-20T00:00:00.000000000Z'": {time.Unix(-8640000000000, 0)},
	}
	for expected, d := range tests {
		if actual, err := d.SurrealString(); assert.NoError(t, err) {
			assert.Equal(t, expected, actual)
		}
	}
}

func TestDatetimeCBOR(t *testing.T) {
	tests := map[string]models.Datetime{
		"d'2024-06-01T12:34:56.780123456Z'": {time.Unix(1717245296, 780123456)},
	}
	for expected, src := range tests {
		data, err := models.CBORFormatter.Marshal(src)
		if assert.NoError(t, err) {
			var dst models.Datetime
			if err := models.CBORFormatter.Unmarshal(data, &dst); assert.NoError(t, err) {
				if actual, err := dst.SurrealString(); assert.NoError(t, err) {
					assert.Equal(t, expected, actual)
				}
			}
		}
	}
}

func TestDatetimeJSON(t *testing.T) {
	tests := map[string]models.Datetime{
		"2024-06-01T12:34:56.780123456Z": {time.Unix(1717245296, 780123456)},
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

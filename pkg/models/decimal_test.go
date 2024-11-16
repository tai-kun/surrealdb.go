package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tai-kun/surrealdb.go/pkg/models"
)

func TestDecimalSurrealString(t *testing.T) {
	tests := map[string]models.Decimal{
		"3.14dec": "3.14",
	}
	for expected, d := range tests {
		if actual, err := d.SurrealString(); assert.NoError(t, err) {
			assert.Equal(t, expected, actual)
		}
	}
}

func TestDecimalCBOR(t *testing.T) {
	tests := map[string]models.Decimal{
		"3.14dec": "3.14",
	}
	for expected, src := range tests {
		data, err := models.CBORFormatter.Marshal(src)
		if assert.NoError(t, err) {
			var dst models.Decimal
			if err := models.CBORFormatter.Unmarshal(data, &dst); assert.NoError(t, err) {
				if actual, err := dst.SurrealString(); assert.NoError(t, err) {
					assert.Equal(t, expected, actual)
				}
			}
		}
	}
}

func TestDecimalJSON(t *testing.T) {
	tests := map[string]models.Decimal{
		"3.14": "3.14",
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

package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tai-kun/surrealdb.go/pkg/models"
)

func TestUUIDSurrealString(t *testing.T) {
	tests := map[string]models.UUID{
		"u'26c80163-3b83-481b-93da-c473947cccbc'": {
			0x26,
			0xc8,
			0x01,
			0x63,
			// -
			0x3b,
			0x83,
			// -
			0x48,
			0x1b,
			// -
			0x93,
			0xda,
			// -
			0xc4,
			0x73,
			0x94,
			0x7c,
			0xcc,
			0xbc,
		},
	}
	for expected, d := range tests {
		if actual, err := d.SurrealString(); assert.NoError(t, err) {
			assert.Equal(t, expected, actual)
		}
	}
}

func TestUUIDCBOR(t *testing.T) {
	tests := map[string]models.UUID{
		"u'26c80163-3b83-481b-93da-c473947cccbc'": {
			0x26,
			0xc8,
			0x01,
			0x63,
			// -
			0x3b,
			0x83,
			// -
			0x48,
			0x1b,
			// -
			0x93,
			0xda,
			// -
			0xc4,
			0x73,
			0x94,
			0x7c,
			0xcc,
			0xbc,
		},
	}
	for expected, src := range tests {
		data, err := models.CBORFormatter.Marshal(src)
		if assert.NoError(t, err) {
			var dst models.UUID
			if err := models.CBORFormatter.Unmarshal(data, &dst); assert.NoError(t, err) {
				if actual, err := dst.SurrealString(); assert.NoError(t, err) {
					assert.Equal(t, expected, actual)
				}
			}
		}
	}
}

func TestUUIDJSON(t *testing.T) {
	tests := map[string]models.UUID{
		"26c80163-3b83-481b-93da-c473947cccbc": {
			0x26,
			0xc8,
			0x01,
			0x63,
			// -
			0x3b,
			0x83,
			// -
			0x48,
			0x1b,
			// -
			0x93,
			0xda,
			// -
			0xc4,
			0x73,
			0x94,
			0x7c,
			0xcc,
			0xbc,
		},
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

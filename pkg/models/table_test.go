package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tai-kun/surrealdb.go/pkg/models"
)

func TestTableSurrealString(t *testing.T) {
	tests := map[string]models.Table{
		"`tai-kun`": "tai-kun",
	}
	for expected, d := range tests {
		if actual, err := d.SurrealString(); assert.NoError(t, err) {
			assert.Equal(t, expected, actual)
		}
	}
}

func TestTableCBOR(t *testing.T) {
	tests := map[string]models.Table{
		"`tai-kun`": "tai-kun",
	}
	for expected, src := range tests {
		data, err := models.CBORFormatter.Marshal(src)
		if assert.NoError(t, err) {
			var dst models.Table
			if err := models.CBORFormatter.Unmarshal(data, &dst); assert.NoError(t, err) {
				if actual, err := dst.SurrealString(); assert.NoError(t, err) {
					assert.Equal(t, expected, actual)
				}
			}
		}
	}
}

func TestTableJSON(t *testing.T) {
	tests := map[string]models.Table{
		"tai-kun": "tai-kun",
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

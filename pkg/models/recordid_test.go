package models_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/tai-kun/surrealdb.go/pkg/models"
)

func TestRecordIDSurrealString(t *testing.T) {
	type value interface {
		SurrealString() (string, error)
	}
	tests := map[string]value{
		`r'⟨tai-kun⟩:1'`:    models.NewRecordID("tai-kun", 1),
		`r'⟨tai-kun⟩:3.14'`: models.NewRecordID("tai-kun", 3.14),
		`r'⟨tai-kun⟩:⟨-1⟩'`: models.NewRecordID("tai-kun", -1),
		`r'city:{"date":"2024-06-01T21:00:00.000000000Z","name":"Tokyo","temp":29.6}'`: models.NewRecordID(
			"city",
			map[string]any{
				"name": "Tokyo",
				"date": models.Datetime{time.Date(2024, 6, 1, 21, 0, 0, 0, time.UTC)},
				"temp": 29.6,
			},
		),
		`r"user:u'26c80163-3b83-481b-93da-c473947cccbc'"`: models.NewRecordID(
			"user",
			models.UUID([16]byte{0x26, 0xc8, 0x01, 0x63, 0x3b, 0x83, 0x48, 0x1b, 0x93, 0xda, 0xc4, 0x73, 0x94, 0x7c, 0xcc, 0xbc}),
		),
	}
	for expected, src := range tests {
		if s, err := src.SurrealString(); assert.NoError(t, err) {
			assert.Equal(t, expected, s)
		}
	}
}

func TestRecordIDCBOR(t *testing.T) {
	tests := map[string]*models.RecordID[int]{
		`r'⟨tai-kun⟩:1'`: models.NewRecordID("tai-kun", 1),
	}
	for expected, src := range tests {
		data, err := models.CBORFormatter.Marshal(src)
		if assert.NoError(t, err) {
			var dst models.RecordID[int]
			if err := models.CBORFormatter.Unmarshal(data, &dst); assert.NoError(t, err) {
				if actual, err := dst.SurrealString(); assert.NoError(t, err) {
					assert.Equal(t, expected, actual)
				}
			}
		}
	}
}

func TestRecordIDJSON(t *testing.T) {
	tests := map[string]*models.RecordID[int]{
		`⟨tai-kun⟩:1`: models.NewRecordID("tai-kun", 1),
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

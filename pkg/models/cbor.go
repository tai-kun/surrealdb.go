package models

import (
	"fmt"
	"reflect"

	"github.com/fxamacker/cbor/v2"

	"github.com/tai-kun/surrealdb.go/pkg/codec"
)

const (
	TagNone          uint64 = 6
	TagTable         uint64 = 7
	TagRecordID      uint64 = 8
	TagDecimal       uint64 = 10
	TagDatetime      uint64 = 12
	TagDuration      uint64 = 14
	TagFuture        uint64 = 15
	TagUUID          uint64 = 37
	TagRange         uint64 = 49
	TagBoundIncluded uint64 = 50
	TagBoundExcluded uint64 = 51
	// TagGeometryPoint        uint64 = 88
	// TagGeometryLine         uint64 = 89
	// TagGeometryPolygon      uint64 = 90
	// TagGeometryMultipoint   uint64 = 91
	// TagGeometryMultiline    uint64 = 92
	// TagGeometryMultipolygon uint64 = 93
	// TagGeometryCollection   uint64 = 94
)

// type Model interface {
// 	MarshalCBOR() ([]byte, error)
// 	UnmarshalCBOR(data []byte) error
// 	MarshalJSON() ([]byte, error)
// 	UnmarshalJSON(data []byte) error
// 	SurrealString() (string, error)
// }

var DefaultModels = map[uint64]any{
	TagNone:          None{},
	TagTable:         Table(""),
	TagRecordID:      RecordID[any]{},
	TagDecimal:       Decimal(""),
	TagDatetime:      Datetime{},
	TagDuration:      Duration(0),
	TagFuture:        Future(""),
	TagUUID:          UUID([16]byte{}),
	TagRange:         Range[any]{},
	TagBoundIncluded: BoundIncluded[any]{},
	TagBoundExcluded: BoundExcluded[any]{},
	// TagGeometryPoint
	// TagGeometryLine
	// TagGeometryPolygon
	// TagGeometryMultipoint
	// TagGeometryMultiline
	// TagGeometryMultipolygon
	// TagGeometryCollection
}

func tagSet() cbor.TagSet {
	tags := cbor.NewTagSet()
	for num, i := range DefaultModels {
		if err := tags.Add(
			cbor.TagOptions{
				EncTag: cbor.EncTagRequired,
				DecTag: cbor.DecTagRequired,
			},
			reflect.TypeOf(i),
			num,
		); err != nil {
			err := fmt.Errorf(
				"surrealdb: models: failed to create CBOR tag set: "+
					"failed to add tagged data item num=%d: %w",
				num, err,
			)
			panic(err)
		}
	}
	return tags
}

var (
	CBORFormatter *codec.CBORFormatter = codec.NewCBORFormatter(tagSet())
	JSONFormatter *codec.JSONFormatter = codec.NewJSONFormatter()
)

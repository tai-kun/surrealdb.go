package surrealdb

import (
	"github.com/tai-kun/surrealdb.go/pkg/codec"
	"github.com/tai-kun/surrealdb.go/pkg/models"
)

var (
	CBORFormatter codec.Formatter = models.CBORFormatter
	JSONFormatter codec.Formatter = models.JSONFormatter
)

// func NewCBORFormatter(override map[uint64]any) codec.Formatter {
// 	tags := cbor.NewTagSet()
// 	for num, i := range models.DefaultModels {
// 		if j, ok := override[num]; ok {
// 			i = j
// 		}
// 		if err := tags.Add(
// 			cbor.TagOptions{
// 				EncTag: cbor.EncTagRequired,
// 				DecTag: cbor.DecTagRequired,
// 			},
// 			reflect.TypeOf(i),
// 			num,
// 		); err != nil {
// 			err := fmt.Errorf(
// 				"surrealdb: failed to create CBOR tag set: "+
// 					"failed to add tagged data item num=%d: %w",
// 				num, err,
// 			)
// 			panic(err)
// 		}
// 	}
// 	return codec.NewCBORFormatter(tags)
// }

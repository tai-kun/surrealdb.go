package codec

import (
	"fmt"
	"io"

	"github.com/fxamacker/cbor/v2"
)

type CBORFormatter struct {
	em cbor.EncMode
	dm cbor.DecMode
}

func NewCBORFormatter(tags cbor.TagSet) *CBORFormatter {
	em, err := cbor.EncOptions{}.EncModeWithTags(tags)
	if err != nil {
		err := fmt.Errorf(
			"surrealdb: codec: failed to create CBOR formatter: "+
				"failed to create cbor.EncMode: %w",
			err,
		)
		panic(err)
	}

	dm, err := cbor.DecOptions{}.DecModeWithTags(tags)
	if err != nil {
		err := fmt.Errorf(
			"surrealdb: codec: failed to create CBOR formatter: "+
				"failed to create cbor.DecMode: %w",
			err,
		)
		panic(err)
	}

	return &CBORFormatter{
		em: em,
		dm: dm,
	}
}

func (cf *CBORFormatter) ContentType() string {
	return "application/cbor"
}

// func (cf *CBORFormatter) WSProtocols() []string {
// 	return []string{"cbor"}
// }

func (cf *CBORFormatter) Marshal(v any) ([]byte, error) {
	return cf.em.Marshal(v)
}

func (cf *CBORFormatter) Unmarshal(data []byte, dst any) error {
	return cf.dm.Unmarshal(data, dst)
}

func (cf *CBORFormatter) NewEncoder(w io.Writer) *cbor.Encoder {
	return cf.em.NewEncoder(w)
}

func (cf *CBORFormatter) NewDecoder(r io.Reader) *cbor.Decoder {
	return cf.dm.NewDecoder(r)
}

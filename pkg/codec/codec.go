package codec

// import (
// 	"io"
// )

// type Encoder interface {
// 	Encode(v any) error
// }

// type Decoder interface {
// 	Decode(v any) error
// }

type Marshaler interface {
	Marshal(v interface{}) ([]byte, error)
	// NewEncoder(w io.Writer) Encoder
}

type Unmarshaler interface {
	Unmarshal(data []byte, dst interface{}) error
	// NewDecoder(r io.Reader) Decoder
}

type Formatter interface {
	Marshaler
	Unmarshaler
	ContentType() string
	// WSProtocols() []string
}
package codec

import (
	"encoding/json"
	"io"
)

type JSONFormatter struct{}

func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

func (jf *JSONFormatter) ContentType() string {
	return "application/json"
}

// func (jf *JSONFormatter) WSProtocols() []string {
// 	return []string{"json"}
// }

func (jf *JSONFormatter) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (jf *JSONFormatter) Unmarshal(data []byte, dst any) error {
	return json.Unmarshal(data, dst)
}

func (jf *JSONFormatter) NewEncoder(w io.Writer) *json.Encoder {
	return json.NewEncoder(w)
}

func (jf *JSONFormatter) NewDecoder(r io.Reader) *json.Decoder {
	return json.NewDecoder(r)
}

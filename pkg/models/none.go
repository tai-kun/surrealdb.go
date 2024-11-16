package models

type None struct{}

func (n *None) MarshalCBOR() ([]byte, error) {
	return []byte{198, 246}, nil
}

func (n *None) UnmarshalCBOR(data []byte) error {
	return nil
}

func (n *None) MarshalJSON() ([]byte, error) {
	return []byte{110, 117, 108, 108}, nil
}

func (n *None) UnmarshalJSON(data []byte) error {
	return nil
}

func (n *None) SurrealString() (string, error) {
	return "NONE", nil
}

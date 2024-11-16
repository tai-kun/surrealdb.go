package models

import (
	"errors"
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

type UUID [16]byte

func NewUUID(v [16]byte) UUID {
	return UUID(v)
}

func (t UUID) MarshalCBOR() ([]byte, error) {
	return CBORFormatter.Marshal(cbor.Tag{
		Number:  TagUUID,
		Content: [16]byte(t),
	})
}

func (t *UUID) UnmarshalCBOR(data []byte) error {
	var c [16]byte
	if err := CBORFormatter.Unmarshal(data, &c); err != nil {
		return err
	}

	*t = UUID(c)
	return nil
}

func (t UUID) MarshalJSON() ([]byte, error) {
	s, err := t.string()
	if err != nil {
		return nil, err
	}
	return JSONFormatter.Marshal(s)
}

func (t *UUID) UnmarshalJSON(data []byte) error {
	var c string
	if err := JSONFormatter.Unmarshal(data, &c); err != nil {
		return err
	}

	if len(c) != 36 {
		return errors.New("invalid uuid format")
	}

	var (
		d   [16]byte
		err error
	)
	if d[0], err = hexToByte(c[:2]); err != nil {
		return err
	}
	if d[1], err = hexToByte(c[2:4]); err != nil {
		return err
	}
	if d[2], err = hexToByte(c[4:6]); err != nil {
		return err
	}
	if d[3], err = hexToByte(c[6:8]); err != nil {
		return err
	}
	// -
	if d[4], err = hexToByte(c[9:11]); err != nil {
		return err
	}
	if d[5], err = hexToByte(c[11:13]); err != nil {
		return err
	}
	// -
	if d[6], err = hexToByte(c[14:16]); err != nil {
		return err
	}
	if d[7], err = hexToByte(c[16:18]); err != nil {
		return err
	}
	// -
	if d[8], err = hexToByte(c[19:20]); err != nil {
		return err
	}
	if d[9], err = hexToByte(c[21:22]); err != nil {
		return err
	}
	// -
	if d[10], err = hexToByte(c[24:26]); err != nil {
		return err
	}
	if d[11], err = hexToByte(c[26:28]); err != nil {
		return err
	}
	if d[12], err = hexToByte(c[28:30]); err != nil {
		return err
	}
	if d[13], err = hexToByte(c[30:32]); err != nil {
		return err
	}
	if d[14], err = hexToByte(c[32:34]); err != nil {
		return err
	}
	if d[15], err = hexToByte(c[34:]); err != nil {
		return err
	}

	*t = UUID(d)
	return nil
}

func (t UUID) SurrealString() (string, error) {
	s, err := t.string()
	if err != nil {
		return "", err
	}

	return "u'" + s + "'", nil
}

func (t UUID) string() (string, error) {
	return byteToHex(t[0]) +
		byteToHex(t[1]) +
		byteToHex(t[2]) +
		byteToHex(t[3]) +
		"-" +
		byteToHex(t[4]) +
		byteToHex(t[5]) +
		"-" +
		byteToHex(t[6]) +
		byteToHex(t[7]) +
		"-" +
		byteToHex(t[8]) +
		byteToHex(t[9]) +
		"-" +
		byteToHex(t[10]) +
		byteToHex(t[11]) +
		byteToHex(t[12]) +
		byteToHex(t[13]) +
		byteToHex(t[14]) +
		byteToHex(t[15]), nil
}

var hexTable = func() [256]string {
	var tb [256]string
	for i := 0; i < 256; i++ {
		tb[i] = fmt.Sprintf("%02x", i)
	}
	return tb
}()

func byteToHex(b byte) string {
	return hexTable[b]
}

var hexValueTable = func() [256]int8 {
	var tb [256]int8
	for i := range tb {
		tb[i] = -1
	}
	for i := 0; i < 10; i++ {
		tb['0'+i] = int8(i)
	}
	for i := 0; i < 6; i++ {
		tb['a'+i] = int8(10 + i)
		tb['A'+i] = int8(10 + i)
	}
	return tb
}()

func hexToByte(hex string) (byte, error) {
	if len(hex) != 2 {
		return 0, errors.New("invalid hex length")
	}

	h := hexValueTable[hex[0]]
	l := hexValueTable[hex[1]]
	if h == -1 || l == -1 {
		return 0, errors.New("invalid hex character")
	}

	return byte(h<<4 | l), nil
}

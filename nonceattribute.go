package stun

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type NonceAttribute struct {
	Nonce string
}

func (h *NonceAttribute) Type() AttributeType {
	return Nonce
}

func (h *NonceAttribute) Encode(msg *Message) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, attributeHeader(Attribute(h), msg))
	err = binary.Write(buf, binary.BigEndian, h.Nonce)

	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (h *NonceAttribute) Decode(data []byte, length uint16, _ *Message) error {
	if uint16(len(data)) < length {
		return errors.New("Truncated Nonce Attribute")
	}
	if length > 763 {
		return errors.New("Nonce Length is too long")
	}
	h.Nonce = string(data[0:length])
	return nil
}

func (h *NonceAttribute) Length(_ *Message) uint16 {
	return uint16(len(h.Nonce))
}

package stun

import (
	"bytes"
	"errors"
)

type UnknownStunAttribute struct {
	ClaimedType AttributeType
	Data        []byte
}

func NewUnknownAttribute() Attribute {
	return Attribute(new(UnknownStunAttribute))
}

func (h *UnknownStunAttribute) Type() AttributeType {
	return h.ClaimedType
}

func (h *UnknownStunAttribute) Encode(msg *Message) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := WriteHeader(buf, Attribute(h), msg)
	buf.Write(h.Data)

	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (h *UnknownStunAttribute) Decode(data []byte, length uint16, _ *Parser) error {
	if uint16(len(data)) < length {
		return errors.New("Truncated Attribute")
	}
	h.Data = data[0:length]
	return nil
}

func (h *UnknownStunAttribute) Length(_ *Message) uint16 {
	return uint16(len(h.Data))
}

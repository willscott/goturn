package stun

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type UnknownStunAttribute struct {
	ClaimedType AttributeType
	Data        []byte
}

func (h *UnknownStunAttribute) Type() AttributeType {
	return h.ClaimedType
}

func (h *UnknownStunAttribute) Encode(_ *Message) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, attributeHeader(Attribute(h)))
	err = binary.Write(buf, binary.BigEndian, h.Data)

	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (h *UnknownStunAttribute) Decode(data []byte, length uint16, _ *Message) error {
	if uint16(len(data)) < length {
		return errors.New("Truncated Attribute")
	}
	h.Data = data[0:length]
	return nil
}

func (h *UnknownStunAttribute) Length() uint16 {
	return uint16(len(h.Data))
}

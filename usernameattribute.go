package stun

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type UsernameAttribute struct {
}

func (h *UsernameAttribute) Type() AttributeType {
	return Username
}

func (h *UsernameAttribute) Encode(msg *Message) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, attributeHeader(Attribute(h), msg))
	err = binary.Write(buf, binary.BigEndian, msg.Credentials.Username)

	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (h *UsernameAttribute) Decode(data []byte, length uint16, msg *Message) error {
	if uint16(len(data)) < length {
		return errors.New("Truncated Username Attribute")
	}
	msg.Credentials.Username = string(data[0:length])
	return nil
}

func (h *UsernameAttribute) Length(msg *Message) uint16 {
	return uint16(len(msg.Credentials.Username))
}

package stun

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type RealmAttribute struct {
}

func (h *RealmAttribute) Type() AttributeType {
	return Realm
}

func (h *RealmAttribute) Encode(msg *Message) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, attributeHeader(Attribute(h), msg))
	err = binary.Write(buf, binary.BigEndian, msg.Credentials.Realm)

	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (h *RealmAttribute) Decode(data []byte, length uint16, msg *Message) error {
	if uint16(len(data)) < length {
		return errors.New("Truncated Realm Attribute")
	}
	if length > 763 {
		return errors.New("Realm Length is too long")
	}
	msg.Credentials.Realm = string(data[0:length])
	return nil
}

func (h *RealmAttribute) Length(msg *Message) uint16 {
	return uint16(len(msg.Credentials.Realm))
}

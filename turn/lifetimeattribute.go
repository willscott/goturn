package turn

import (
	"bytes"
	"encoding/binary"
	"errors"
  "github.com/willscott/goturn"
)

type LifetimeAttribute struct {
  Lifetime uint32
}

func (h *LifetimeAttribute) Type() stun.AttributeType {
	return Lifetime
}

func (h *LifetimeAttribute) Encode(_ *stun.Message) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, attributeHeader(stun.Attribute(h), msg))
	err = binary.Write(buf, binary.BigEndian, h.Lifetime)

	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (h *LifetimeAttribute) Decode(data []byte, length uint16, _ *stun.Message) error {
	if uint16(len(data)) < length {
		return errors.New("Truncated Username Attribute")
	}
  h.Lifetime = uint32(data[0:4])
	return nil
}

func (h *LifetimeAttribute) Length(_ *stun.Message) uint16 {
	return 4
}

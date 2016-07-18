package turn

import (
	"bytes"
	"encoding/binary"
	"errors"
  "github.com/willscott/goturn"
)

type ChannelNumberAttribute struct {
  ChannelNumber uint16
}

func (h *ChannelNumberAttribute) Type() stun.AttributeType {
	return ChannelNumber
}

func (h *ChannelNumberAttribute) Encode(msg *stun.Message) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, attributeHeader(stun.Attribute(h), msg))
	err = binary.Write(buf, binary.BigEndian, h.ChannelNumber)
  err = binary.Write(buf, binary.BigEndian, uint16(0))

	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (h *ChannelNumberAttribute) Decode(data []byte, length uint16, _ *stun.Message) error {
	if length != 4 || uint16(len(data)) < length {
		return errors.New("Truncated ChannelNumber Attribute")
	}
  h.ChannelNumber = binary.BigEndian.Uint16(data[0:2])
	return nil
}

func (h *ChannelNumberAttribute) Length(_ *stun.Message) uint16 {
	return 4
}

package stun

import (
  "bytes"
  "encoding/binary"
  "errors"
)

type MessageIntegrityAttribute struct {
  Hash []byte
}

func (h *MessageIntegrityAttribute) Type() (StunAttributeType) {
  return MessageIntegrity
}

func (h *MessageIntegrityAttribute) Encode(_ *StunMessage) ([]byte, error) {
  buf := new(bytes.Buffer)
  err := binary.Write(buf, binary.BigEndian, attributeHeader(StunAttribute(h)))
  err = binary.Write(buf, binary.BigEndian, h.Hash)

  if err != nil {
    return nil, err
  }
  return buf.Bytes(), nil
}

func (h *MessageIntegrityAttribute) Decode(data []byte, length uint16, _ *Header) (error) {
  if length != 20 || len(data) < 20 {
    return errors.New("Truncated MessageIntegrity Attribute")
  }
  h.Hash = data[0:20]
  return nil
}

func (h *MessageIntegrityAttribute) Length() (uint16) {
  return 20
}

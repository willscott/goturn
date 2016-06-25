package turn

import (
  "bytes"
  "encoding/binary"
  "errors"
  "log"
)

type UnknownStunAttribute struct {
  ClaimedType StunAttributeType
  Data []byte
}

func (h UnknownStunAttribute) Type() (StunAttributeType) {
  return h.ClaimedType
}

func (h UnknownStunAttribute) Encode() ([]byte, error) {
  buf := new(bytes.Buffer)
  err := binary.Write(buf, binary.BigEndian, attributeHeader(StunAttribute(h)))
  err = binary.Write(buf, binary.BigEndian, h.Data)

  if err != nil {
    return nil, err
  }
  return buf.Bytes(), nil
}

func (h UnknownStunAttribute) Decode(data []byte, length uint16) (error) {
  if uint16(len(data)) < length {
    return errors.New("Truncated Attribute")
  }
  h.Data = data[0:length]
  log.Print("len is", len(h.Data))
  return nil
}

func (h UnknownStunAttribute) Length() (uint16) {
  return uint16(len(h.Data))
}

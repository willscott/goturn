package turn

import (
  "bytes"
  "encoding/binary"
  "errors"
)

type UsernameAttribute struct {
  Username string
}

func (h *UsernameAttribute) Type() (StunAttributeType) {
  return Username
}

func (h *UsernameAttribute) Encode() ([]byte, error) {
  buf := new(bytes.Buffer)
  err := binary.Write(buf, binary.BigEndian, attributeHeader(StunAttribute(h)))
  err = binary.Write(buf, binary.BigEndian, h.Username)

  if err != nil {
    return nil, err
  }
  return buf.Bytes(), nil
}

func (h *UsernameAttribute) Decode(data []byte, length uint16) (error) {
  if uint16(len(data)) < length {
    return errors.New("Truncated Username Attribute")
  }
  h.Username = string(data[0:length])
  return nil
}

func (h *UsernameAttribute) Length() (uint16) {
  return uint16(len(h.Username))
}

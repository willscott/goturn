package turn

import (
  "bytes"
  "encoding/binary"
  "errors"
)

type FingerprintAttribute struct {
  CRC []byte
}

func (h *FingerprintAttribute) Type() (StunAttributeType) {
  return Fingerprint
}

func (h *FingerprintAttribute) Encode() ([]byte, error) {
  buf := new(bytes.Buffer)
  err := binary.Write(buf, binary.BigEndian, attributeHeader(StunAttribute(h)))
  err = binary.Write(buf, binary.BigEndian, h.CRC)

  if err != nil {
    return nil, err
  }
  return buf.Bytes(), nil
}

func (h *FingerprintAttribute) Decode(data []byte, length uint16) (error) {
  if length != 4 || len(data) < 4 {
    return errors.New("Truncated Fingerprint Attribute")
  }
  h.CRC = data[0:4]
  return nil
}

func (h *FingerprintAttribute) Length() (uint16) {
  return 4
}

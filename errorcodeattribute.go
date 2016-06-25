package turn

import (
  "bytes"
  "encoding/binary"
  "errors"
)

type ErrorCodeAttribute struct {
  Class uint8
  Number uint8
  Phrase string
}

func (h ErrorCodeAttribute) Type() (StunAttributeType) {
  return ErrorCode
}

func (h ErrorCodeAttribute) Encode() ([]byte, error) {
  buf := new(bytes.Buffer)
  err := binary.Write(buf, binary.BigEndian, attributeHeader(StunAttribute(h)))
  err = binary.Write(buf, binary.BigEndian, uint16(0))
  err = binary.Write(buf, binary.BigEndian, h.Class)
  err = binary.Write(buf, binary.BigEndian, h.Number)
  err = binary.Write(buf, binary.BigEndian, h.Phrase)

  if err != nil {
    return nil, err
  }
  return buf.Bytes(), nil
}

func (h ErrorCodeAttribute) Decode(data []byte, length uint16) (error) {
  if len(data) < 4 || len(data) > 65535 || uint16(len(data)) < length {
    return errors.New("Truncated Error Code Attribute")
  }
  if uint8(data[0]) != 0 || uint8(data[1]) != 0 {
    return errors.New("Invalid reserved bytes in Error Code Attribute")
  }
  h.Class = uint8(data[2])
  if h.Class < 3 || h.Class > 6 {
    return errors.New("Invalid Error Code Class")
  }
  h.Number = uint8(data[3])
  if h.Number > 99 {
    return errors.New("Invlaid Error Code Number")
  }
  h.Phrase = string(data[4:length])
  return nil
}

func (h ErrorCodeAttribute) Length() (uint16) {
  return uint16(4 + len(h.Phrase))
}

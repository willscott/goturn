package turn

import (
  "bytes"
  "encoding/binary"
  "errors"
)

const (
  magicCookie uint32 = 0x2112A442
)

type StunClass uint16
const (
  StunRequest StunClass = iota
  StunIndication
  StunResponse
  StunError
)

type StunType uint16
const (
  StunBinding StunType = 1 + iota
)

type StunHeader struct {
  Class   StunClass
  Type    StunType
  Length  uint16
  Id      []byte
}

type StunAttributeType uint16
const (
  MappedAddress StunAttributeType = 0x1
  Username = 0x6
  MessageIntegrity = 0x8
  ErrorCode = 0x9
  UnknownAttributes = 0xA
  Realm = 0x14
  Nonce = 0x15
  XorMappedAddress = 0x20

  // comprehension-optional attributes
  Software = 0x8022
  AlternateServer = 0x8023
  Fingerprint = 0x8028
)

type StunAttribute struct {
  Type    StunAttributeType
  Length  uint16
  Value   []byte
}

func (h *StunHeader) Encode() ([]byte, error) {
  var classEnc uint16 = 0
  buf := new(bytes.Buffer)

  hType := uint16(h.Type)
  hClass := uint16(h.Class)

  //bits 0-3 are low bits of type
  classEnc |= hType & 15
  //bit 4 is low bit of class
  classEnc |= (hClass & 1) << 4
  //bits 5-7 are bits 4-6 of type
  classEnc |= ((hType >> 4) & 7) << 5
  //bit 8 is high bit of class
  classEnc |= (hClass & 2) << 7
  //bits 9-13 are high bits of type
  classEnc |= ((hType >> 7) & 31) << 9

  err := binary.Write(buf, binary.BigEndian, classEnc)
  err = binary.Write(buf, binary.BigEndian, h.Length)
  err = binary.Write(buf, binary.BigEndian, magicCookie)
  err = binary.Write(buf, binary.BigEndian, h.Id)

  if len(h.Id) != 12 {
    return nil, errors.New("Unsupported Transaction ID Length")
  }

  if err != nil {
    return nil, err
  }
  return buf.Bytes(), nil
}

func (h *StunHeader) Decode(data []byte) (error) {
  if len(data) < 20 {
    return errors.New("Header Length Too Short")
  }

  classEnc := binary.BigEndian.Uint16(data)
  stunClass := StunClass(((classEnc & 4) >> 3) + ((classEnc & 8) >> 6))
  stunType := StunType(classEnc & 15 + ((classEnc >> 5) & 7) << 4 + ((classEnc >> 9) & 31) << 7)

  if classEnc >> 14 != 0 {
    return errors.New("First 2 bits are not 0")
  }

  if binary.BigEndian.Uint32(data[4:]) != magicCookie {
    return errors.New("Bad Magic Cookie")
  }

  if binary.BigEndian.Uint16(data[2:]) & 3 != 0 {
    return errors.New("Message Length is not a multiple of 4")
  }

  h.Type = stunType
  h.Class = stunClass
  h.Length = binary.BigEndian.Uint16(data[2:])
  h.Id = data[8:20]

  return nil
}

package turn

import (
  "bytes"
  "encoding/binary"
  "errors"
  "fmt"
  "net"
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

type StunAttribute interface {
  Type()          StunAttributeType
  Encode()        ([]byte, error)
  Decode([]byte)  error
  Length()        uint16
}

type StunMessage struct {
  Header StunHeader
  Attributes []StunAttribute
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

func DecodeStunAttribute(data []byte) (*StunAttribute, error) {
  attributeType := binary.BigEndian.Uint16(data)
  length := binary.BigEndian.Uint16(data[2:])
  var result StunAttribute
  switch StunAttributeType(attributeType) {
  case MappedAddress:
    result = new(MappedAddressAttribute)
  }
  err := result.Decode(data[4:])
  if err != nil {
    return nil, err
  } else if result.Length() != length {
    return nil, errors.New(fmt.Sprintf("Incorrect Length Specified for %T", result))
  }
  return &result, nil
}

func attributeHeader(a StunAttribute) (uint32) {
  attributeType := uint16(a.Type())
  return (uint32(attributeType) << 16) + uint32(a.Length())
}

type MappedAddressAttribute struct {
  Family  uint16
  Port    uint16
  Address net.IP
}

func (h MappedAddressAttribute) Type() (StunAttributeType) {
  return MappedAddress
}

func (h MappedAddressAttribute) Encode() ([]byte, error) {
  buf := new(bytes.Buffer)
  err := binary.Write(buf, binary.BigEndian, attributeHeader(StunAttribute(h)))
  err = binary.Write(buf, binary.BigEndian, h.Family)
  err = binary.Write(buf, binary.BigEndian, h.Port)
  err = binary.Write(buf, binary.BigEndian, h.Address)

  if err != nil {
    return nil, err
  }
  return buf.Bytes(), nil
}

func (h MappedAddressAttribute) Decode(data []byte) (error) {
  if data[0] != 0 && data[1] != 1 && data[0] != 2 {
    return errors.New("Incorrect Mapped Address Family.")
  }
  h.Family = uint16(data[1])
  if (h.Family == 1 && len(data) < 8) || (h.Family == 2 && len(data) < 20) {
    return errors.New("Mapped Address Attribute unexpectedly Truncated.")
  }
  h.Port = uint16(data[2]) << 8 + uint16(data[3])
  if h.Family == 1 {
    h.Address = data[4:8]
  } else {
    h.Address = data[4:20]
  }
  return nil
}

func (h MappedAddressAttribute) Length() (uint16) {
  if h.Family == 1 {
    return 8
  } else {
    return 20
  }
}

func Parse(data []byte) (*StunMessage, error) {
  message := new(StunMessage)
  message.Attributes = []StunAttribute{}
  if err := message.Header.Decode(data); err != nil {
    return nil, err
  }
  data = data[20:]
  if len(data) != int(message.Header.Length) {
    return nil, errors.New("Message has incorrect Length")
  }
  for len(data) > 0 {
    attribute := new(StunAttribute)
    if err := (*attribute).Decode(data); err != nil {
      return nil, err
    }
    message.Attributes = append(message.Attributes, *attribute)
    len := int((*attribute).Length() + 3 / 4)
    data = data[len:]
  }
  return message, nil
}

func (m *StunMessage) Serialize() ([]byte, error) {
  body := []byte{}

  // Calculate length.
  m.Header.Length = uint16(len(body))
  head, err := m.Header.Encode()
  if err != nil {
    return nil, err
  }
  data := append(head, body...)
  return data, nil
}

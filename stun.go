package turn

import (
  "bytes"
  "crypto/rand"
  "encoding/binary"
  "errors"
  "fmt"
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
  Id      [12]byte
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
  Decode([]byte, uint16)  error
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
  copy(h.Id[:], data[8:20])

  return nil
}

func DecodeStunAttribute(data []byte) (*StunAttribute, error) {
  attributeType := binary.BigEndian.Uint16(data)
  length := binary.BigEndian.Uint16(data[2:])
  var result StunAttribute
  switch StunAttributeType(attributeType) {
  case MappedAddress:
    result = new(MappedAddressAttribute)
  case Username:
    result = new(UsernameAttribute)
  default:
    unknownAttr := new(UnknownStunAttribute)
    unknownAttr.ClaimedType = StunAttributeType(attributeType)
    result = unknownAttr
  }
  err := result.Decode(data[4:], length)
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
    attribute, err := DecodeStunAttribute(data)
    if err != nil {
      return nil, err
    }
    message.Attributes = append(message.Attributes, *attribute)
    // 4 byte header and rounded up to next multiple of 4
    len := 4 * int(((*attribute).Length() + 7) / 4)
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

//Convienence functions for making commonly used data structures.
func MakeStunRequest(header *StunHeader) (error) {
  header.Class = StunRequest
  header.Type = StunBinding
  _, err := rand.Read(header.Id[:])
  return err
}

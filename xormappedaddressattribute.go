package stun

import (
  "bytes"
  "encoding/binary"
  "errors"
  "net"
)

type XorMappedAddressAttribute struct {
  Family  uint16
  Port    uint16
  Address net.IP
}

func (h *XorMappedAddressAttribute) Type() (AttributeType) {
  return XorMappedAddress
}

func (h *XorMappedAddressAttribute) Encode(msg *Message) ([]byte, error) {
  buf := new(bytes.Buffer)
  err := binary.Write(buf, binary.BigEndian, attributeHeader(Attribute(h)))
  err = binary.Write(buf, binary.BigEndian, h.Family)
  xport := h.Port ^ uint16(magicCookie >> 16)
  err = binary.Write(buf, binary.BigEndian, xport)

  var xoraddress []byte
  if h.Family == 1 {
    xoraddress = make([]byte, 4)
    binary.BigEndian.PutUint32(xoraddress, magicCookie)
  } else {
    xoraddress = make([]byte, 16)
    binary.BigEndian.PutUint32(xoraddress, magicCookie)
    copy(xoraddress[4:16], msg.Header.Id[:])
  }
  for i, _ := range xoraddress {
    xoraddress[i] ^= h.Address[i]
  }
  err = binary.Write(buf, binary.BigEndian, xoraddress)

  if err != nil {
    return nil, err
  }
  return buf.Bytes(), nil
}

func (h *XorMappedAddressAttribute) Decode(data []byte, _ uint16, header *Header) (error) {
  if data[0] != 0 && data[1] != 1 && data[0] != 2 {
    return errors.New("Incorrect Mapped Address Family.")
  }
  h.Family = uint16(data[1])
  if (h.Family == 1 && len(data) < 8) || (h.Family == 2 && len(data) < 20) {
    return errors.New("Mapped Address Attribute unexpectedly Truncated.")
  }
  h.Port = uint16(data[2]) << 8 + uint16(data[3])
  // X-port is XOR'ed with the 16 most significant bits of the magic Cookie
  h.Port ^= uint16(magicCookie >> 16)

  var xoraddress []byte
  if h.Family == 1 {
    xoraddress = make([]byte, 4)
    binary.BigEndian.PutUint32(xoraddress, magicCookie)
    h.Address = data[4:8]
  } else {
    xoraddress = make([]byte, 16)
    binary.BigEndian.PutUint32(xoraddress, magicCookie)
    copy(xoraddress[4:16], header.Id[:])
    h.Address = data[4:20]
  }
  for i, _ := range xoraddress {
    h.Address[i] ^= xoraddress[i]
  }
  return nil
}

func (h *XorMappedAddressAttribute) Length() (uint16) {
  if h.Family == 1 {
    return 8
  } else {
    return 20
  }
}

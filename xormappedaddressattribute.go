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

func (h *XorMappedAddressAttribute) Type() (StunAttributeType) {
  return XorMappedAddress
}

func (h *XorMappedAddressAttribute) Encode() ([]byte, error) {
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

func (h *XorMappedAddressAttribute) Decode(data []byte, _ uint16) (error) {
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

func (h *XorMappedAddressAttribute) Length() (uint16) {
  if h.Family == 1 {
    return 8
  } else {
    return 20
  }
}

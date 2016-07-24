package turn

import (
  "bytes"
	common "github.com/willscott/goturn/common"
  "github.com/willscott/goturn/stun"
	"net"
)

const (
	XorPeerAddress common.AttributeType = 0x12
)

type XorPeerAddressAttribute struct {
	Family  uint16
	Port    uint16
	Address net.IP
}

func NewXorPeerAddressAttribute() common.Attribute {
	return common.Attribute(new(XorPeerAddressAttribute))
}

func (h *XorPeerAddressAttribute) Type() common.AttributeType {
	return XorPeerAddress
}

func (h *XorPeerAddressAttribute) Encode(msg *common.Message) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := common.WriteHeader(buf, common.Attribute(h), msg); err != nil {
		return nil, err
	}
  mapped := stun.XorMappedAddressAttribute(*h)
  bytes, err := stun.XorAddressData(&mapped, msg)
  if err != nil {
  	return nil, err
  }
  buf.Write(bytes)
  return buf.Bytes(), nil
}

func (h *XorPeerAddressAttribute) Decode(data []byte, length uint16, p *common.Parser) error {
  mapped := stun.XorMappedAddressAttribute(*h)
  if err := mapped.Decode(data, length, p); err != nil {
    return err
  }
  h.Family = mapped.Family
  h.Port = mapped.Port
  h.Address = mapped.Address
  return nil
}

func (h *XorPeerAddressAttribute) Length(_ *common.Message) uint16 {
	if h.Family == 1 {
		return 8
	} else {
		return 20
	}
}

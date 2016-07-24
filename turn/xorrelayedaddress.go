package turn

import (
	common "github.com/willscott/goturn/common"
  "github.com/willscott/goturn/stun"
	"net"
)

const (
	XorRelayedAddress common.AttributeType = 0x16
)

type XorRelayedAddressAttribute struct {
	Family  uint16
	Port    uint16
	Address net.IP
}

func NewXorRelayedAddressAttribute() common.Attribute {
	return common.Attribute(new(XorRelayedAddressAttribute))
}

func (h *XorRelayedAddressAttribute) Type() common.AttributeType {
	return XorRelayedAddress
}

func (h *XorRelayedAddressAttribute) Encode(msg *common.Message) ([]byte, error) {
  mapped := stun.XorMappedAddressAttribute(*h)
  return mapped.Encode(msg)
}

func (h *XorRelayedAddressAttribute) Decode(data []byte, length uint16, p *common.Parser) error {
  mapped := stun.XorMappedAddressAttribute(*h)
  if err := mapped.Decode(data, length, p); err != nil {
    return err
  }
  h.Family = mapped.Family
  h.Port = mapped.Port
  h.Address = mapped.Address
  return nil
}

func (h *XorRelayedAddressAttribute) Length(_ *common.Message) uint16 {
	if h.Family == 1 {
		return 8
	} else {
		return 20
	}
}

package stun

import (
	"bytes"
	"errors"
	"github.com/willscott/goturn/common"
)

const (
	Nonce stun.AttributeType = 0x15
)

type NonceAttribute struct {
	Nonce []byte
}

func NewNonceAttribute() stun.Attribute {
	return stun.Attribute(new(NonceAttribute))
}

func (h *NonceAttribute) Type() stun.AttributeType {
	return Nonce
}

func (h *NonceAttribute) Encode(msg *stun.Message) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := stun.WriteHeader(buf, stun.Attribute(h), msg)
	buf.Write(h.Nonce)

	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (h *NonceAttribute) Decode(data []byte, length uint16, _ *stun.Parser) error {
	if uint16(len(data)) < length {
		return errors.New("Truncated Nonce Attribute")
	}
	if length > 763 {
		return errors.New("Nonce Length is too long")
	}
	h.Nonce = data[0:length]
	return nil
}

func (h *NonceAttribute) Length(_ *stun.Message) uint16 {
	return uint16(len(h.Nonce))
}

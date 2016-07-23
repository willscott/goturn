package stun

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/willscott/goturn/common"
)

const (
	ErrorCode stun.AttributeType = 0x9
)

type ErrorCodeAttribute struct {
	Class  uint8
	Number uint8
	Phrase string
}

func NewErrorCodeAttribute() stun.Attribute {
	return stun.Attribute(new(ErrorCodeAttribute))
}

func (h *ErrorCodeAttribute) Type() stun.AttributeType {
	return ErrorCode
}

func (h *ErrorCodeAttribute) Encode(msg *stun.Message) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := stun.WriteHeader(buf, stun.Attribute(h), msg)
	err = binary.Write(buf, binary.BigEndian, uint16(0))
	err = binary.Write(buf, binary.BigEndian, h.Class)
	err = binary.Write(buf, binary.BigEndian, h.Number)
	err = binary.Write(buf, binary.BigEndian, h.Phrase)

	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (h *ErrorCodeAttribute) Decode(data []byte, length uint16, _ *stun.Message) error {
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
		return errors.New("Invalid Error Code Number")
	}
	h.Phrase = string(data[4:length])
	return nil
}

func (h *ErrorCodeAttribute) Length(_ *stun.Message) uint16 {
	return uint16(4 + len(h.Phrase))
}

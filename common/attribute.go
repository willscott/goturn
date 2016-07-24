package stun

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type AttributeType uint16

type Attribute interface {
	Type() AttributeType
	Encode(*Message) ([]byte, error)
	Decode([]byte, uint16, *Parser) error
	Length(*Message) uint16
}

type AttributeSet map[AttributeType]func() Attribute

func WriteHeader(buf *bytes.Buffer, a Attribute, msg *Message) error {
	attributeType := uint16(a.Type())
	header := (uint32(attributeType) << 16) + uint32(a.Length(msg))
	return binary.Write(buf, binary.BigEndian, header)
}

func DecodeAttribute(data []byte, attrs AttributeSet, parser *Parser) (*Attribute, error) {
	attributeType := binary.BigEndian.Uint16(data)
	length := binary.BigEndian.Uint16(data[2:])
	attrMaker, ok := attrs[AttributeType(attributeType)]
	if !ok {
		attrMaker = NewUnknownAttribute
	}
	result := attrMaker()

	err := result.Decode(data[4:], length, parser)
	if err != nil {
		return nil, err
	} else if result.Length(parser.Message) != length {
		return nil, errors.New(fmt.Sprintf("Incorrect Length Specified for %T", result))
	}
	return &result, nil
}

package stun

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

type HeaderType uint16

const (
	BindingRequest       HeaderType = 0x0001
	SharedSecretRequest             = 0x0002
	BindingResponse                 = 0x0101
	SharedSecretResponse            = 0x0102
	BindingError                    = 0x0111
	SharedSecretError               = 0x0112
)

type Header struct {
	Type   HeaderType
	Length uint16
	Id     [12]byte
}

type Credentials struct {
	Username string
	Realm    string
	Password string
}

func (h Header) String() string {
	return fmt.Sprintf("%T #%x [%db]", h.Type, h.Id, h.Length)
}

type AttributeType uint16

const (
	MappedAddress     AttributeType = 0x1
	Username                        = 0x6
	MessageIntegrity                = 0x8
	ErrorCode                       = 0x9
	UnknownAttributes               = 0xA
	Realm                           = 0x14
	Nonce                           = 0x15
	XorMappedAddress                = 0x20

	// comprehension-optional attributes
	Software        = 0x8022
	AlternateServer = 0x8023
	Fingerprint     = 0x8028
)

type Attribute interface {
	Type() AttributeType
	Encode(*Message) ([]byte, error)
	Decode([]byte, uint16, *Message) error
	Length(*Message) uint16
}

type Message struct {
	Header
	Credentials
	Attributes []Attribute
}

func (h *Header) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.BigEndian, h.Type)
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

func (h *Header) Decode(data []byte) error {
	if len(data) < 20 {
		return errors.New("Header Length Too Short")
	}

	// Correctness checks.
	if binary.BigEndian.Uint16(data[0:])>>14 != 0 {
		return errors.New("First 2 bits are not 0")
	}

	if binary.BigEndian.Uint32(data[4:]) != magicCookie {
		return errors.New("Bad Magic Cookie")
	}

	if binary.BigEndian.Uint16(data[2:])&3 != 0 {
		return errors.New("Message Length is not a multiple of 4")
	}

	h.Type = HeaderType(binary.BigEndian.Uint16(data[0:]))
	h.Length = binary.BigEndian.Uint16(data[2:])
	copy(h.Id[:], data[8:20])

	return nil
}

func DecodeAttribute(data []byte, msg *Message) (*Attribute, error) {
	attributeType := binary.BigEndian.Uint16(data)
	length := binary.BigEndian.Uint16(data[2:])
	var result Attribute
	switch AttributeType(attributeType) {
	case ErrorCode:
		result = new(ErrorCodeAttribute)
	case Fingerprint:
		result = new(FingerprintAttribute)
	case MappedAddress:
		result = new(MappedAddressAttribute)
	case MessageIntegrity:
		result = new(MessageIntegrityAttribute)
	case Nonce:
		result = new(NonceAttribute)
	case Realm:
		result = new(RealmAttribute)
	case Username:
		result = new(UsernameAttribute)
	case XorMappedAddress:
		result = new(XorMappedAddressAttribute)
	default:
		unknownAttr := new(UnknownStunAttribute)
		unknownAttr.ClaimedType = AttributeType(attributeType)
		result = unknownAttr
	}
	err := result.Decode(data[4:], length, msg)
	if err != nil {
		return nil, err
	} else if result.Length(msg) != length {
		return nil, errors.New(fmt.Sprintf("Incorrect Length Specified for %T", result))
	}
	return &result, nil
}

func attributeHeader(a Attribute, msg *Message) uint32 {
	attributeType := uint16(a.Type())
	return (uint32(attributeType) << 16) + uint32(a.Length(msg))
}

func Parse(data []byte, credentials Credentials) (*Message, error) {
	message := new(Message)
	message.Credentials = credentials
	message.Attributes = []Attribute{}
	if err := message.Header.Decode(data); err != nil {
		return nil, err
	}
	data = data[20:]
	if len(data) != int(message.Header.Length) {
		return nil, errors.New("Message has incorrect Length")
	}
	for len(data) > 0 {
		attribute, err := DecodeAttribute(data, message)
		if err != nil {
			return nil, err
		}
		message.Attributes = append(message.Attributes, *attribute)
		// 4 byte header and rounded up to next multiple of 4
		len := 4 * int(((*attribute).Length(message)+7)/4)
		data = data[len:]
	}
	return message, nil
}

func (m *Message) Serialize() ([]byte, error) {
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
func NewBindingRequest() (*Message, error) {
	message := Message{
		Header: Header{
			Type: BindingRequest,
		},
	}
	_, err := rand.Read(message.Header.Id[:])
	return &message, err
}

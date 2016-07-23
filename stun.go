package goturn

import (
	"crypto/rand"
	"errors"
	common "github.com/willscott/goturn/common"
	"github.com/willscott/goturn/stun"
)

const (
	BindingRequest       common.HeaderType = 0x0001
	SharedSecretRequest                    = 0x0002
	BindingResponse                        = 0x0101
	SharedSecretResponse                   = 0x0102
	BindingError                           = 0x0111
	SharedSecretError                      = 0x0112
)

//Deprecated. Should live in individual stun attribute implementations.
const (
	AlternateServer common.AttributeType = 0x8023
)

func Parse(data []byte, credentials common.Credentials, attrs common.AttributeSet) (*common.Message, error) {
	message := new(common.Message)
	message.Credentials = credentials
	message.Attributes = []common.Attribute{}
	if err := message.Header.Decode(data); err != nil {
		return nil, err
	}
	data = data[20:]
	if len(data) != int(message.Header.Length) {
		return nil, errors.New("Message has incorrect Length")
	}
	for len(data) > 0 {
		attribute, err := common.DecodeAttribute(data, attrs, message)
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

func ParseStun(data []byte) (*common.Message, error) {
	return Parse(data, common.Credentials{}, stun.StunAttributes)
}

//Convienence functions for making commonly used data structures.
func NewBindingRequest() (*common.Message, error) {
	message := common.Message{
		Header: common.Header{
			Type: BindingRequest,
		},
	}
	_, err := rand.Read(message.Header.Id[:])
	return &message, err
}

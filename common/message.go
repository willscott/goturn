package stun

import (
	"errors"
)

type Message struct {
	Header
	Credentials
	Attributes []Attribute
}

func (m *Message) Serialize() ([]byte, error) {
	var bodylength uint16
	for _, att := range m.Attributes {
    len := uint16(4 * int((att.Length(m)+7)/4))
		bodylength += len
	}

	body := make([]byte, bodylength)
	bodylength = 0

	for _, att := range m.Attributes {
		attLen := uint16(4 * int((att.Length(m)+7)/4))
		if attBody, err := att.Encode(m); err != nil {
			return nil, err
		} else {
			copy(body[bodylength:bodylength+attLen], attBody[:])
		}
		bodylength += attLen
	}

	// Calculate length.
	m.Header.Length = uint16(len(body))
	head, err := m.Header.Encode()
	if err != nil {
		return nil, err
	}
	data := append(head, body...)
	return data, nil
}

func (m *Message) GetAttribute(typ AttributeType) *Attribute {
  for _, att := range m.Attributes {
    if att.Type() == typ {
      return &att
    }
  }
  return nil
}

func Parse(data []byte, credentials Credentials, attrs AttributeSet) (*Message, error) {
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
		attribute, err := DecodeAttribute(data, attrs, message)
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

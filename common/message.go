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

type Parser struct {
  *Message
  *Credentials
  AttributeSet
  Data []byte
  Offset uint16
}

func Parse(data []byte, credentials *Credentials, attrs AttributeSet) (*Message, error) {
  parser := Parser{new(Message), credentials, attrs, data, 0}
  err := parser.Parse()
  if err != nil {
    return nil, err
  }
  return parser.Message, nil
}

func (p *Parser) Parse() error {
	if p.Credentials != nil {
		p.Message.Credentials = *p.Credentials
	}
	p.Message.Attributes = []Attribute{}
	if err := p.Message.Header.Decode(p.Data); err != nil {
		return err
	}
  data := p.Data[20:]
  p.Offset = 20
	if len(data) != int(p.Message.Header.Length) {
		return errors.New("Message has incorrect Length")
	}
	for len(data) > 0 {
		attribute, err := DecodeAttribute(data, p.AttributeSet, p)
		if err != nil {
			return err
		}
		p.Message.Attributes = append(p.Message.Attributes, *attribute)
		// 4 byte header and rounded up to next multiple of 4
		len := 4 * int(((*attribute).Length(p.Message)+7)/4)
    p.Offset += uint16(len)
		data = data[len:]
	}
	return nil
}

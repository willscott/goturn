package stun

import (
	"fmt"
)

type Credentials struct {
	Username string
	Realm    string
	Password string
}

func (c Credentials) String() string {
	return fmt.Sprintf("%s:%s@%s", c.Username, c.Password, c.Realm)
}

type Message struct {
	Header
	Credentials
	Attributes []Attribute
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

func (m *Message) GetAttribute(typ AttributeType) *Attribute {
  for _, att := range m.Attributes {
    if att.Type() == typ {
      return &att
    }
  }
  return nil
}

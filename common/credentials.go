package stun

import (
	"fmt"
)

type Credentials struct {
	Nonce    []byte
	Username string
	Realm    string
	Password string
}

func (c Credentials) String() string {
	return fmt.Sprintf("%s:%s@%s [nonce %s]", c.Username, c.Password, c.Realm, c.Nonce)
}

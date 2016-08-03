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

func (c *Credentials) ForNewConnection() *Credentials {
	creds := new(Credentials)
	creds.Username = c.Username
	creds.Realm = c.Realm
	creds.Password = c.Password
	return creds
}

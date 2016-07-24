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

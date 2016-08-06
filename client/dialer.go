package client

import (
	"github.com/willscott/goturn/common"
	"net"
	"time"
)

// A TURN Dialer is a dialer implementation which dials through a TURN relay.
type TurnDialer struct {
	// The connection to the Relay
	StunClient

	// Fields to implement a net.Dialer
	Timeout       time.Duration
	Deadline      time.Time
	LocalAddr     net.Addr
	DualStack     bool
	FallbackDelay time.Duration
	KeepAlive     time.Duration
	Cancel        <-chan struct{}
}

func NewDialer(credentials *stun.Credentials, control net.Conn) (d *TurnDialer, err error) {
  d = new(TurnDialer)
	d.StunClient.Conn = control

	addr, err := d.StunClient.Allocate(credentials)
	if err != nil {
		return nil, err
	}
  d.LocalAddr = addr

  //TODO: functional cancel channel.

	return d, nil
}

func (d *TurnDialer) Dial(network, addr string) (c net.Conn, err error) {
	endpoint := stun.NewAddressFromString(network, addr)
	if err := d.StunClient.RequestPermission(endpoint); err != nil {
		return nil, err
	}

	c, err = d.StunClient.Connect(endpoint)
	if err != nil {
		return nil, err
	}
	return c, nil
}

//Utility function to avoid need for clients to directly import common.
func LongtermCredentials(username, password string) stun.Credentials {
  return stun.Credentials{Username: username, Password: password}
}

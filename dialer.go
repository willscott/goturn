package turn

import (
  "net"
  "time"
)

// A TURN Dialer is a dialer implementation which dials through a TURN relay.
type Dialer struct {
  // The connection to the TURN Relay
  net.Conn

  // Fields to implement a net.Dialer
  Timeout time.Duration
  Deadline time.Time
  LocalAddr net.Addr
  DualStack bool
  FallbackDelay time.Duration
  KeepAlive time.Duration
  Cancel <-chan struct{}
}

type turn struct {
  user, password string
  network, addr string
  forward net.Dialer
}

// Create a TURN Dialer for a given relay.
/*
func TURNDialer(network, addr string) (Dialer, err error) {
  dialer := Dialer{}
  dialer.Conn, err = DialRelay(network, addr, dialer.Timeout)
  if err != nil {
    return nil, err
  }
  return dialer, nil
}

// Dial a relay.
func DialRelay(network, addr string, timeout time.Duration) (conn *Conn, err error) {
  conn = new(net.Conn)
  conn.Conn, err = net.DialTimeout(network, address, timeout)
  if err != nil {
    return nil, err
  }
  return conn, nil
}

func (d *Dialer) Dial(network, addr string) (c net.Conn, err error) {
  switch network {
  case "udp", "udp6", "udp4":
  default:
    return nil, errors.New("proxy: no support for TURN proxy connections of type " + network)
  }
}
*/

package client

import (
  "net"
  "github.com/willscott/goturn/common"
  "github.com/willscott/goturn"
)

// The Client maintains state on a connection with a stun/turn server

type STUNClient struct {
  // The connection handles requests to the server, and can be either UDP,
  // TCP, or TCP over TLS.
  net.Conn
}

func (s *StunClient) readStunPacket() (*stun.Message, error) {
  // listen for response
  s.Conn.SetReadDeadline(time.Now().Add(1000 * time.Millisecond))
  b := make([]byte, 2048)
  n, err := c.Read(b)
  if err != nil || n == 0 || n > 2048 {
    return nil, err
  }
  msg, err := goturn.ParseStun(b[:n])
	if err != nil {
    return nil, err
	}
  return msg, nil
}

func (s *StunClient) Bind() (net.Addr, error) {
  // construct request message
  packet, err := goturn.NewBindingRequest()
  if err != nil {
    return nil, err
  }

  message, err := packet.Serialize()
  if err != nil {
    return nil, err
  }

  // send message
  if _, err = s.Conn.Write(message); err != nil {
    return nil, err
  }

  response,err := s.readStunPacket()
  if err != nil {
    return nil, err
  }

  if response.Header.Type != goturn.BindingResponse {
    return nil, errors.New("Unexpected response type.")
	}
  response.GetAttribute()
}

func (s *StunClient) Allocate(stun.Credentials) error {

}

func (s *StunClient) RequestPermission(net.Addr) error {

}

func (s *StunClient)

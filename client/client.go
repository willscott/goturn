package client

import (
	"bufio"
	"errors"
	"github.com/willscott/goturn"
	"github.com/willscott/goturn/common"
	stunattrs "github.com/willscott/goturn/stun"
	turnattrs "github.com/willscott/goturn/turn"
	"net"
	"time"
)

// The Client maintains state on a connection with a stun/turn server

type StunClient struct {
	// The connection handles requests to the server, and can be either UDP,
	// TCP, or TCP over TLS.
	net.Conn

	// A buffered reader helps us read from the connection.
	reader *bufio.Reader

	// The dialer for making new connections.
	net.Dialer

	// The credentials used for authenticating communication with the server.
	*stun.Credentials

	// Timeout until the active connection expires.
	Timeout time.Duration

	// Time until the next message must be received.
	Deadline time.Time
}

// Create a new connection to the same remote endpoint, sharing credentials
// with the current connection
func (s *StunClient) deriveConnection() (*StunClient, error) {
	other := new(StunClient)
	other.Dialer = s.Dialer
	other.Credentials = s.Credentials.ForNewConnection()
	other.Timeout = s.Timeout

	conn, err := s.Dialer.Dial(s.Conn.RemoteAddr().Network(), s.Conn.RemoteAddr().String())
	if err != nil {
		return nil, err
	}
	other.Conn = conn
	return other, nil
}

// Send a message from the client.
func (s *StunClient) send(packet *stun.Message, err error) error {
	if err != nil {
		return err
	}
	packet.Credentials = *s.Credentials

	message, err := packet.Serialize()
	if err != nil {
		return err
	}

	// send message
	// TODO: handle not all of message being written.
	if _, err = s.Conn.Write(message); err != nil {
		return err
	}
	return nil
}

// Read the next packet off of the connection abstracted by the client.
// Returns either the next message, or an error if the next set of bytes
// do not represent a valid message.
func (s *StunClient) readStunPacket() (*stun.Message, error) {
	// Set up timeouts for reading.
	if s.reader == nil {
		s.reader = bufio.NewReader(s.Conn)
	}
	if s.Timeout > 0 {
		// Initial Deadline
		if s.Deadline.IsZero() {
			s.Deadline = time.Now().Add(s.Timeout)
		}
		s.Conn.SetReadDeadline(s.Deadline)
	}

	// Start by reading the header to learn the length of the packet.
	h, err := s.reader.Peek(20)
	if err != nil {
		return nil, err
	}
	header := stun.Header{}
	if err = header.Decode(h); err != nil {
		return nil, err
	}
	if header.Length == 0 {
		return goturn.ParseTurn(h, s.Credentials)
	}
	if header.Length > 2048 {
		return nil, errors.New("Packet length too long.")
	}
	buffer := make([]byte, 20+header.Length)
	n, err := s.reader.Read(buffer)
	if err != nil || uint16(n) != 20+header.Length {
		return nil, err
	}

	if s.Timeout > 0 {
		s.Deadline = time.Now().Add(s.Timeout)
	}

	return goturn.ParseTurn(buffer, s.Credentials)
}

// Request a Stun Binding to learn the Internet-visible address of the current
// connection.
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

	response, err := s.readStunPacket()
	if err != nil {
		return nil, err
	}

	if response.Header.Type != goturn.BindingResponse {
		return nil, errors.New("Unexpected response type.")
	}
	attr := response.GetAttribute(stunattrs.MappedAddress)
	port := uint16(0)
	address := net.IP{}

	if attr != nil {
		addr := (*attr).(*stunattrs.MappedAddressAttribute)
		port = addr.Port
		address = addr.Address
	} else {
		attr = response.GetAttribute(stunattrs.XorMappedAddress)
		if attr == nil {
			return nil, errors.New("No Mapped Address provided.")
		}
		addr := (*attr).(*stunattrs.XorMappedAddressAttribute)
		port = addr.Port
		address = addr.Address
	}
	return stun.NewAddress(s.Conn.RemoteAddr().Network(), address, port), nil
}

func (s *StunClient) allocateUnauthenticated() error {
	// make a simple allocation message
	creds := s.Credentials
	s.Credentials = &stun.Credentials{}
	if err := s.send(goturn.NewAllocateRequest(s.Conn.RemoteAddr().Network(), false)); err != nil {
		return err
	}
	s.Credentials = creds

	response, err := s.readStunPacket()
	if err != nil {
		return err
	}

	if response.Credentials.Nonce != nil {
		s.Credentials.Nonce = response.Credentials.Nonce
	}
	if len(response.Credentials.Realm) > 0 {
		s.Credentials.Realm = response.Credentials.Realm
	}
	msgerr := stunattrs.GetError(response)
	if msgerr.Error() > 0 && msgerr.Error() != 401 {
		return errors.New("Initial Connection failed " + msgerr.String())
	}
	return nil
}

// Request to connect to a Turn server. The Turn protocol uses the term
// allocation to refer to an authenticated connection with the server.
// Returns the bound address.
func (s *StunClient) Allocate(c *stun.Credentials) (net.Addr, error) {
	s.Credentials = c

	if s.Credentials.Nonce == nil {
		if err := s.allocateUnauthenticated(); err != nil {
			return nil, err
		}
	}
	if err := s.send(goturn.NewAllocateRequest(s.Conn.RemoteAddr().Network(), true)); err != nil {
		return nil, err
	}
	response, err := s.readStunPacket()
	if err != nil {
		return nil, err
	}

	if response.Header.Type != goturn.AllocateResponse {
		msgerr := stunattrs.GetError(response)
		if msgerr.Error() == 442 {
			// TODO: bad transport; retry w/ other protocol.
		}
		return nil, errors.New("Connection failed: " + msgerr.String())
	}

	relayAddr := response.GetAttribute(turnattrs.XorRelayedAddress)
	relayAddress := (*relayAddr).(*turnattrs.XorRelayedAddressAttribute)

	return stun.NewAddress(s.Conn.RemoteAddr().Network(), relayAddress.Address, relayAddress.Port), nil
}

// Request permission to relay data with a remote address. The Client should
// already have an authenticated connection with the server, using Allocate,
// for this request to succeed.
func (s *StunClient) RequestPermission(with net.Addr) error {
	addr := stun.Address{with}
	if err := s.send(goturn.NewPermissionRequest(addr.HostPart())); err != nil {
		return err
	}
	response, err := s.readStunPacket()
	if err != nil {
		return err
	}

	if response.Header.Type != goturn.CreatePermissionResponse {
		return errors.New("Connection failed: " + stunattrs.GetError(response).String())
	}
	return nil
}

//Assumes that there is already an allocation for the client.
func (s *StunClient) Connect(to net.Addr) (net.Conn, error) {
	if err := s.send(goturn.NewConnectRequest(to)); err != nil {
		return nil, err
	}
	response, err := s.readStunPacket()
	if err != nil {
		return nil, err
	}

	if response.Header.Type != goturn.ConnectResponse {
		return nil, errors.New("Connection failed: " + stunattrs.GetError(response).String())
	}

	// extract Connection-id
	connID := response.GetAttribute(turnattrs.ConnectionId)
	if connID == nil {
		return nil, errors.New("No Connection ID provided.")
	}
	connectionID := (*connID).(*turnattrs.ConnectionIdAttribute).ConnectionId

	// create the data connection.
	conn, err := s.deriveConnection()
	if err != nil {
		return nil, err
	}

	if err := conn.send(goturn.NewConnectionBindRequest(connectionID)); err != nil {
		conn.Conn.Close()
		return nil, err
	}

	response, err = conn.readStunPacket()
	if err != nil {
		conn.Conn.Close()
		return nil, err
	}

	// Need to get nonce for the new connection first.
	if response.Credentials.Nonce != nil {
		conn.Credentials.Nonce = response.Credentials.Nonce
	}

	if err := conn.send(goturn.NewConnectionBindRequest(connectionID)); err != nil {
		conn.Conn.Close()
		return nil, err
	}

	response, err = conn.readStunPacket()
	if err != nil {
		conn.Conn.Close()
		return nil, err
	}

	if response.Header.Type != goturn.ConnectionBindResponse {
		conn.Conn.Close()
		return nil, errors.New("Connection failed: " + stunattrs.GetError(response).String())
	}

	return conn.Conn, nil
}

package goturn

import (
	"crypto/rand"

	common "github.com/willscott/goturn/common"
	"github.com/willscott/goturn/stun"
	"github.com/willscott/goturn/turn"

	"net"
)

const (
	AllocateRequest          common.HeaderType = 0x0003
	RefreshRequest                             = 0x0004
	CreatePermissionRequest                    = 0x0008
	ChannelBindRequest                         = 0x0009
	ConnectRequest                             = 0x000a
	ConnectionBindRequest                      = 0x000b
	SendIndication                             = 0x0016
	DataIndication                             = 0x0017
	ConnectionAttemptIndication                = 0x001c
	AllocateResponse                           = 0x0103
	RefreshResponse                            = 0x0104
	CreatePermissionResponse                   = 0x0108
	ChannelBindResponse                        = 0x0109
	ConnectResponse                            = 0x010a
	ConnectionBindResponse                     = 0x010b
	AllocateError                              = 0x0113
	RefreshError                               = 0x0114
	CreatePermissionError                      = 0x0118
	ChannelBindError                           = 0x0119
	ConnectError                               = 0x011a
	ConnectionBindError                        = 0x011b
)

//Deprecated. Should live in individual turn attribute implementations.
const (
	EvenPort           common.AttributeType = 0x18
	DontFragment                            = 0x1A
	ReservationToken                        = 0x22
)

func ParseTurn(data []byte, credentials common.Credentials) (*common.Message, error) {
	return common.Parse(data, credentials, turn.AttributeSet())
}

func NewAllocateRequest(inResponseTo *common.Message) (*common.Message, error) {
	message := common.Message{
		Header: common.Header{
			Type: AllocateRequest,
		},
	}
	_, err := rand.Read(message.Header.Id[:])

	//Include a RequestedTransportAttribute = UDP
	if inResponseTo == nil {
		message.Attributes = []common.Attribute{&turn.RequestedTransportAttribute{17}}
	} else {
		message.Credentials = inResponseTo.Credentials
		message.Attributes = []common.Attribute{&turn.RequestedTransportAttribute{17},
			&stun.NonceAttribute{},
			&stun.UsernameAttribute{},
			&stun.RealmAttribute{},
			&stun.MessageIntegrityAttribute{},
			&stun.FingerprintAttribute{}}
	}

	return &message, err
}

func NewPermissionRequest(credentials common.Credentials, to net.IP) (*common.Message, error) {
	message := common.Message{
		Header: common.Header{
			Type: CreatePermissionRequest,
		},
	}
	_, err := rand.Read(message.Header.Id[:])

	message.Credentials = credentials
	family := uint16(1)
	if to.To4() == nil {
		family = 2
	}
	message.Attributes = []common.Attribute{&stun.NonceAttribute{},
		&turn.XorPeerAddressAttribute{family, 0, to},
		&stun.UsernameAttribute{},
		&stun.RealmAttribute{},
		&stun.MessageIntegrityAttribute{},
		&stun.FingerprintAttribute{}}

	return &message, err
}

func NewSendIndication(host net.IP, port uint16, data []byte) (*common.Message, error) {
	message := common.Message{
		Header: common.Header{
			Type: SendIndication,
		},
	}
	_, err := rand.Read(message.Header.Id[:])
	family := uint16(1)
	if host.To4() == nil {
		family = 2
	}
	message.Attributes = []common.Attribute{&turn.XorPeerAddressAttribute{family, port, host},
		&turn.DataAttribute{data}}

	return &message, err
}

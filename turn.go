package goturn

import (
	"crypto/rand"

	common "github.com/willscott/goturn/common"
	"github.com/willscott/goturn/turn"
)

const (
	AllocateRequest          common.HeaderType = 0x0003
	RefreshRequest                             = 0x0004
	CreatePermissionRequest                    = 0x0008
	ChannelBindRequest                         = 0x0009
	SendIndication                             = 0x0016
	DataIndication                             = 0x0017
	AllocateResponse                           = 0x0103
	RefreshResponse                            = 0x0104
	CreatePermissionResponse                   = 0x0108
	ChannelBindResponse                        = 0x0109
	AllocateError                              = 0x0113
	RefreshError                               = 0x0114
	CreatePermissionError                      = 0x0118
	ChannelBindError                           = 0x0119
)

const (
	ChannelNumber      common.AttributeType = 0xC
	Lifetime                                = 0xD
	XorPeerAddress                          = 0x12
	Data                                    = 0x13
	XorRelayedAddress                       = 0x16
	EvenPort                                = 0x18
	RequestedTransport                      = 0x19
	DontFragment                            = 0x1A
	ReservationToken                        = 0x22
)

func ParseTurn(data []byte, credentials common.Credentials) (*common.Message, error) {
	return Parse(data, credentials, turn.TurnAttributes)
}

func NewAllocateRequest() (*common.Message, error) {
	message := common.Message{
		Header: common.Header{
			Type: AllocateRequest,
		},
	}
	_, err := rand.Read(message.Header.Id[:])

	//Include a RequestedTransportAttribute = UDP
	message.Attributes = []common.Attribute{&turn.RequestedTransportAttribute{17}}

	return &message, err
}

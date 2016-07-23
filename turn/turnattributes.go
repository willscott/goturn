package turn

import (
	common "github.com/willscott/goturn/common"
	//"github.com/willscott/goturn/stun"
)

var (
	TurnAttributes = common.AttributeSet{
		ChannelNumber:      NewChannelNumberAttribute,
		Lifetime:           NewLifetimeAttribute,
		RequestedTransport: NewRequestedTransportAttribute,
	}
)

package stun

import (
	"github.com/willscott/goturn/common"
)

var (
	StunAttributes = stun.AttributeSet{
		ErrorCode:         NewErrorCodeAttribute,
		Fingerprint:       NewFingerprintAttribute,
		MappedAddress:     NewMappedAddressAttribute,
		MessageIntegrity:  NewMessageIntegrityAttribute,
		Nonce:             NewNonceAttribute,
		Realm:             NewRealmAttribute,
		Software:          NewSoftwareAttribute,
		UnknownAttributes: NewUnknownAttributesAttribute,
		Username:          NewUsernameAttribute,
		XorMappedAddress:  NewXorMappedAddressAttribute,
	}
)

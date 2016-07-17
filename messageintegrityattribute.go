package stun

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/binary"
	"errors"
)

type MessageIntegrityAttribute struct {
}

func (h *MessageIntegrityAttribute) Type() AttributeType {
	return MessageIntegrity
}

func makeKey(msg *Message) []byte {
	var key []byte
	if len(msg.Credentials.Username) > 0 {
		sum := md5.Sum([]byte(msg.Credentials.Username + ":" + msg.Credentials.Realm + ":" + msg.Credentials.Password))
		copy(key[:], sum[0:16])
		return key
	} else if len(msg.Credentials.Password) > 0 {
		return []byte(msg.Credentials.Password)
	} else {
		return nil
	}
}

func (h *MessageIntegrityAttribute) Encode(msg *Message) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, attributeHeader(Attribute(h), msg)); err != nil {
		return nil, err
	}

	key := makeKey(msg)
	if key == nil {
		return nil, errors.New("Cannot sign request without credentials.")
	}

	// Calculate partial message
	var partialMsg Message
	partialMsg.Header = msg.Header
	copy(partialMsg.Attributes, msg.Attributes)

	// Remove either 1 (msg integrity) or 2 (fingerprint and msg integrity) attributes
	partialMsg.Attributes = partialMsg.Attributes[0 : len(partialMsg.Attributes)-1]
	if partialMsg.Attributes[len(partialMsg.Attributes)-1].Type() == MessageIntegrity {
		partialMsg.Attributes = partialMsg.Attributes[0 : len(partialMsg.Attributes)-1]
	}
	// Add a new attribute w/ same length as msg integrity
	dummy := UnknownStunAttribute{MessageIntegrity, make([]byte, 20)}
	partialMsg.Attributes = append(partialMsg.Attributes, &dummy)
	// calcualte the byte string
	msgBytes, err := partialMsg.Serialize()
	if err != nil {
		return nil, err
	}

	//hmac all but the dummy attribute
	mac := hmac.New(sha1.New, key)
	mac.Write(msgBytes[0 : len(msgBytes)-24])
	hash := mac.Sum(nil)

	err = binary.Write(buf, binary.BigEndian, hash)

	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (h *MessageIntegrityAttribute) Decode(data []byte, length uint16, msg *Message) error {
	if length != 20 || len(data) < 20 {
		return errors.New("Truncated MessageIntegrity Attribute")
	}

	key := makeKey(msg)
	// Calculate partial message
	var partialMsg Message
	partialMsg.Header = msg.Header
	copy(partialMsg.Attributes, msg.Attributes)

	// Add a new attribute w/ same length as fingerprint
	dummy := UnknownStunAttribute{MessageIntegrity, make([]byte, 20)}
	partialMsg.Attributes = append(partialMsg.Attributes, &dummy)
	// calcualte the byte string
	msgBytes, err := partialMsg.Serialize()
	if err != nil {
		return err
	}

	//hmac all but the dummy attribute
	mac := hmac.New(sha1.New, key)
	mac.Write(msgBytes[0 : len(msgBytes)-24])
	hash := mac.Sum(nil)

	if !bytes.Equal(hash, data[0:20]) {
		return errors.New("Invalid Message Integrity value..")
	}

	return nil
}

func (h *MessageIntegrityAttribute) Length(_ *Message) uint16 {
	return 20
}

package stun

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hash/crc32"
)

const (
	crcXOR uint32 = 0x5354554e
)

type FingerprintAttribute struct {
	CRC uint32
}

func (h *FingerprintAttribute) Type() AttributeType {
	return Fingerprint
}

func (h *FingerprintAttribute) Encode(msg *Message) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, attributeHeader(Attribute(h), msg))

	// Calculate partial message
	var partialMsg Message
	partialMsg.Header = msg.Header
	copy(partialMsg.Attributes, msg.Attributes)

	// Fingerprint must be last attribute.
	partialMsg.Attributes = partialMsg.Attributes[0 : len(partialMsg.Attributes)-1]

	// Add a new attribute w/ same length as msg integrity
	dummy := UnknownStunAttribute{Fingerprint, make([]byte, 4)}
	partialMsg.Attributes = append(partialMsg.Attributes, &dummy)
	// calcualte the byte string
	msgBytes, err := partialMsg.Serialize()
	if err != nil {
		return nil, err
	}

	crc := crc32.ChecksumIEEE(msgBytes[0:len(msgBytes)-8]) ^ crcXOR
	err = binary.Write(buf, binary.BigEndian, crc)

	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (h *FingerprintAttribute) Decode(data []byte, length uint16, msg *Message) error {
	if length != 4 || len(data) < 4 {
		return errors.New("Truncated Fingerprint Attribute")
	}
	h.CRC = binary.BigEndian.Uint32(data[0:4])

	// Calculate partial message
	var partialMsg Message
	partialMsg.Header = msg.Header
	copy(partialMsg.Attributes, msg.Attributes)

	// Add a new attribute w/ same length as fingerprint
	dummy := UnknownStunAttribute{Fingerprint, make([]byte, 4)}
	partialMsg.Attributes = append(partialMsg.Attributes, &dummy)
	// calcualte the byte string
	msgBytes, err := partialMsg.Serialize()
	if err != nil {
		return err
	}

	crc := crc32.ChecksumIEEE(msgBytes[0:len(msgBytes)-8]) ^ crcXOR

	if crc != h.CRC {
		return errors.New("Invalid Fingerprint value.")
	}

	return nil
}

func (h *FingerprintAttribute) Length(_ *Message) uint16 {
	return 4
}

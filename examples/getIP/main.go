package main

import (
	"flag"
	"github.com/willscott/goturn"
	"github.com/willscott/goturn/stun"
	"log"
	"net"
	"time"
)

var server = flag.String("server", "stun.l.google.com:19302", "Remote Stun Server")

func parseResponse(datagram []byte) (address net.IP, port uint16) {
	msg, err := goturn.ParseStun(datagram)
	if err != nil {
		log.Fatal("Could not parse response as a STUN message:", err)
	}

	if msg.Header.Type != goturn.BindingResponse {
		log.Fatal("Response message is not a STUN response.", msg.Header)
	}

	for _, attr := range msg.Attributes {
		if attr.Type() == stun.MappedAddress {
			addr := attr.(*stun.MappedAddressAttribute)
			return addr.Address, addr.Port
		} else if attr.Type() == stun.XorMappedAddress {
			addr := attr.(*stun.XorMappedAddressAttribute)
			return addr.Address, addr.Port
		}
	}

	log.Fatal("No MappedAddress in STUN response.")
	return nil, 0
}

func main() {
	flag.Parse()

	// dial
	c, err := net.Dial("udp", *server)
	if err != nil {
		log.Fatal("Could open UDP Connection:", err)
	}
	defer c.Close()

	// construct request message
	packet, err := goturn.NewBindingRequest()
	if err != nil {
		log.Fatal("Failed to generate request packet:", err)
	}

	message, err := packet.Serialize()
	if err != nil {
		log.Fatal("Failed to serialize packet: ", err)
	}

	// send message
	_, err = c.Write(message)
	if err != nil {
		log.Fatal("Failed to send message: ", err)
	}

	// listen for response
	c.SetReadDeadline(time.Now().Add(1000 * time.Millisecond))
	b := make([]byte, 2048)
	n, err := c.Read(b)
	if err != nil || n == 0 || n > 2048 {
		log.Fatal("Failed to read response: ", err)
	}

	address, port := parseResponse(b[:n])
	log.Printf("%s:%d", address, port)
}

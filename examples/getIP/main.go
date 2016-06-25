package main

import (
  "flag"
  "github.com/willscott/goturn"
  "log"
  "net"
  "time"
)

var server = flag.String("server", "stun.l.google.com:19302", "Remote Stun Server")

func parseResponse(datagram []byte) {
  msg, err := turn.Parse(datagram)
  if err != nil {
    log.Fatal("Could not parse response:", err)
  }

  if msg.Header.Class != turn.StunResponse {
    log.Fatal("Response message is not a STUN response.")
  }

  for _, attr := range msg.Attributes {
    if attr.Type() == turn.MappedAddress {
      addr := attr.(*turn.MappedAddressAttribute)
      log.Printf("%s:%d", addr.Address, addr.Port)
      return
    }
  }

  log.Fatal("No MappedAddress in STUN response.")
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
  var packet turn.StunMessage
  if err = turn.MakeStunRequest(&packet.Header); err != nil {
    log.Fatal("Failed to construct header: ", err)
  }
  packet.Attributes = make([]turn.StunAttribute,0)

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
  parseResponse(b[:n])
}

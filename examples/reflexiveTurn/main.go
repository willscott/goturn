package main

import (
  "flag"
  "github.com/willscott/goturn/turn"
  "log"
  "net"
  "time"
)

var credentialURL = flag.String("credentials", "https://computeengineondemand.appspot.com/turn?username=prober&key=4080218913", "credential URL")

func main() {
  flag.Parse()

  // get & parse credentials

  // dial
  c, err := net.Dial("udp", *server)
  if err != nil {
    log.Fatal("Could open UDP Connection:", err)
  }
  defer c.Close()

  // construct request message
  packet,err := turn.NewAllocateRequest()
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

  //address, port := parseResponse(b[:n])
  //log.Printf("%s:%d", address, port)
}

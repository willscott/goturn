package main

import (
  "encoding/json"
  "flag"
  "github.com/willscott/goturn"
  "io/ioutil"
  "log"
  "net"
  "net/http"
  "net/url"
  "time"
)

var credentialURL = flag.String("credentials", "https://computeengineondemand.appspot.com/turn?username=prober&key=4080218913", "credential URL")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Uris []string `json:"uris"`
}

func main() {
  flag.Parse()

  // get & parse credentials
  resp, err := http.Get(*credentialURL)
  if err != nil {
    log.Fatal("Could not fetch URL:", err)
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Fatal("Could not read response:", err)
  }

  var creds Credentials
  if err := json.Unmarshal(body, &creds); err != nil {
    log.Fatal("Could not parse response:", err)
  }

  // Use the first one.
  server, err := url.Parse(creds.Uris[0])
  if err != nil {
    log.Fatal("Invalid server URI:", err)
  }

  log.Printf("Requesting turn to %s", server.Opaque)

  // dial
  c, err := net.Dial("udp", server.Opaque)
  if err != nil {
    log.Fatal("Could open UDP Connection:", err)
  }
  defer c.Close()

  // construct request message
  packet,err := goturn.NewAllocateRequest()
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

Go TURN
=======

This is a library providing a Go `Dialer` interface compatible with the golang
[proxy](https://golang.org/x/net/proxy) package which connects through a
[TURN](https://tools.ietf.org/html/rfc5766) relay.

Installation
------------

```golang
go get github.com/willscott/goturn
```

Full Example
------------

```golang
import (
  "fmt"
  "net"

  "github.com/willscott/goturn"
  "github.com/miekg/dns"
)

dialTurn, err:= goturn.TURNDialer(goturn.UDP, "turn:127.0.0.1")
if err != nil {
  fatalf("Failed to obtain dialer: %v\n", err)
}

udpConn, err := dialTurn.Dial("udp", "8.8.8.8:53")
if err != nil {
  fatalf("Failed to dial: %v\n", err)
}

co := &dns.Conn{Conn: udpConn}
m := new(dns.Msg)
m.SetQuestion("google.com.", dns.TypeA)
co.WriteMsg(m)
in, _ := co.ReadMsg()
co.Close()
if t, ok := in.Answer[0].(*dns.A); ok {
  fmt.Printf("Successfully Resolved google.com to %v through turn\n", t)
}
```

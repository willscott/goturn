Go TURN
=======

This is a library providing a Go interface compatible with the golang
[proxy](https://golang.org/x/net/proxy) package which connects through a
[TURN](https://tools.ietf.org/html/rfc5766) relay.

This package provides parsing and encoding support for [STUN](https://tools.ietf.org/html/rfc5389)
and [TURN](https://tools.ietf.org/html/rfc5766) protocols.

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
  "ioutil"
  "net"
  "net/http"

  "github.com/willscott/goturn/client"
)

// Connect to the stun/turn server
conn, err := net.Dial("tcp", "127.0.0.1:19302")
if err != nil {
  log.Fatal("Could open TCP Connection:", err)
}
defer c.Close()

credentials := client.LongtermCredentials("username", "password")
dialer, err := client.NewDialer(credentials, conn
if err != nil {
  fatalf("Failed to obtain dialer: %v\n", err)
}

httpClient := &http.Client{Transport: &http.Transport{Dial: dialer.Dial}}
httpResp, err := httpClient.Get("http://www.google.com/")
if err != nil {
  log.Fatal("Failed to get webpage", err)
}
defer httpResp.Body.Close()
httpBody, err := ioutil.ReadAll(httpResp.Body)
if err != nil {
  log.Fatal("Failed to read response", err)
}
log.Printf("Received Webpage Body is: %s", string(httpBody))
```

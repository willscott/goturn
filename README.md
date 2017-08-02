Go TURN
=======
[![GoDoc](https://godoc.org/github.com/willscott/goturn?status.svg)](https://godoc.org/github.com/willscott/goturn)

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
package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"

	arg "github.com/alexflint/go-arg"
	"github.com/willscott/goturn/client"
)

func main() {
	var args struct {
		Address  string
		Username string
		Password string
	}
	arg.MustParse(&args)

	// Connect to the stun/turn server
	conn, err := net.Dial("tcp", args.Address)
	if err != nil {
		log.Fatal("error dialing TURN server: ", err)
	}
	defer conn.Close()

	credentials := client.LongtermCredentials(args.Username, args.Password)
	dialer, err := client.NewDialer(&credentials, conn)
	if err != nil {
		log.Fatal("failed to obtain dialer: ", err)
	}

	httpClient := &http.Client{Transport: &http.Transport{Dial: dialer.Dial}}
	httpResp, err := httpClient.Get("http://www.google.com/")
	if err != nil {
		log.Fatal("error performing http request: ", err)
	}
	defer httpResp.Body.Close()

	httpBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		log.Fatal("error reading http response: ", err)
	}
	log.Printf("received %d bytes", len(httpBody))
}
```

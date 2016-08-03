package main

import (
	"encoding/json"
	"flag"
	"github.com/willscott/goturn/client"
  "github.com/willscott/goturn/common"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
)

var server = flag.String("url", "http://myip.info", "URL to fetch")
var credentialURL = flag.String("credentials", "https://computeengineondemand.appspot.com/turn?username=prober&key=4080218913", "credential URL")

type Credentials struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Uris     []string `json:"uris"`
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

  raddr, err := net.ResolveTCPAddr("tcp", server.Opaque)
  if err != nil {
    log.Fatal("Could resolve remote address:", err)
  }
  log.Printf("Negotiating with %s", server.Opaque)

	// dial
	c, err := net.Dial("tcp", raddr.String())
	if err != nil {
		log.Fatal("Could open TCP Connection:", err)
	}
	defer c.Close()

	client := client.StunClient{Conn: c}
  credentials := stun.Credentials{Username: creds.Username, Password: creds.Password}
  if _, err = client.Allocate(&credentials); err != nil {
    log.Fatal("Could not authenticate with server: ", err)
  }
  log.Printf("Authenticated with server.")
}

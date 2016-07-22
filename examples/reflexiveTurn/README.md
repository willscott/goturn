TURN loop-back example
=======================

This example uses TURN to send a packet back to itself. It demonstrates creating
a reservation, and transmitting data through the TURN proxy. This example
assumes an AppRTC-like system for credential distribution, where a URL contains
a JSON-encoded set of credentials.

Usage
-----

```
go run main.go [--credentials https://example.com/credentials.json]
```

Status on reservation, establishment, and transmission of data will be
logged to standard out.

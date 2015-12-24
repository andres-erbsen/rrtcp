This repository contains implementation and testing code for Round Robin TCP, a new transport layer network protocol which uses a round-robin pool of TCP connections to facilitate better bidirectional audio over TCP.

## Dependencies

 The implementation and testing code is written in the `Go` programming language. The TOR tests and implementation require `Go 1.4.1` or greater to run.

## Protocols

Go implementations of UDP, TCP, and Round Robin TCP are located in the `fnet` folder. They use the Go `net` library which uses kernel implementations of these protocols.

## Tests

`*-clock-station` runs simple peer-to-peer tests for each network protocol.
To run the tests yourself, you must run two instances of a given clockstation, one as the listener and one as the dialer, both connected to the same address.

For example, to run the Round Robin TCP tests:

```
cd rrtcp-clock-station
go build
./rrtcp-clock-station -address "localhost:1111" -l
./rrtcp-clock-station -address "localhost:1111"
```

This listener has arbitrarily been chosen as the sender. Every 50 ms, it will send a packet and print the application send time of the packet, using the code in `clockstation`.

The dialer receiver will print the send time and corresponding receive time for each packet it receives, using the code in `clockprinter`.

## Tor
`udp-over-tor` contains code that uses a Go Tor implementation to send packets using Round Robin TCP over TOR.

`udp-over-tor-test` contains code that performs tests analogous to the ones above, bug using Tor. It includes an argument `-deterministic` to deterministcally select a TOR path given a seed, for the purpose of conducting symmetrical tests.

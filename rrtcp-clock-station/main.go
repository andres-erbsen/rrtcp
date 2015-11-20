package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/andres-erbsen/rrtcp/clockprinter"
	"github.com/andres-erbsen/rrtcp/clockstation"
	"github.com/andres-erbsen/rrtcp/fnet"
)

var addr = flag.String("address", "", "address to connect to or listen at")
var listen = flag.Bool("l", false, "bind to the specified address and listen (default: connect)")
var frameSize = flag.Int("s", 1024, "frame size")
var numStreams = flag.Int("n", 5, "number of streams")

func main() {
	flag.Parse()
	if len(flag.Args()) != 0 || *listen == false && *addr == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *listen {
		err := listener(frameSize, numStreams, addr)
		if err != nil {
			os.Exit(2) // TODO: More appropriate per-case error number
		}
	} else {
		err := dialer(frameSize, numStreams, addr)
		if err != nil {
			os.Exit(2)
		}
	}
}

func listener(frameSize *int, numStreams *int, addr *string) error {
	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "net.Listen(%q): %s\n", *addr, err.Error())
		return err
	}
	rrs := fnet.NewStream(*frameSize)

	defer rrs.Stop()
	for i := 0; i < *numStreams; i++ {
		c, err := ln.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ln.Accept(): %s\n", err.Error())
			return err
		}
		rrs.AddStream(c)
	}
	fc := fnet.FrameConn(rrs)
	err = clockstation.Run(fc, time.Tick(50*time.Millisecond))
	if err != nil {
		fmt.Fprintf(os.Stderr, "clockstation.Run: %s\n", err.Error())
		return err
	}
	return nil
}

func dialer(frameSize *int, numStreams *int, addr *string) error {
	rrs := fnet.NewStream(*frameSize)
	defer rrs.Stop()
	for i := 0; i < *numStreams; i++ {
		c, err := net.Dial("tcp", *addr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "net.Dial(%q): %s\n", *addr, err.Error())
			return err
		}
		rrs.AddStream(c)
	}
	fc := fnet.FrameConn(rrs)
	err := clockprinter.Run(fc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "clockprinter.Run: %s\n", err.Error())
		return err
	}
	return nil
}

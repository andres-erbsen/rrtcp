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
var numStreams = flag.Int("n", 2, "number of streams")

func main() {
	flag.Parse()
	if len(flag.Args()) != 0 || *listen == false && *addr == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *listen {
		ln, err := net.Listen("tcp", *addr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "net.Listen(%q): %s\n", *addr, err.Error())
			os.Exit(2)
		}
		rrs := fnet.NewStream(*frameSize)
		for i := 0; i < *numStreams; i++ {
			c, err := ln.Accept()
			if err != nil {
				fmt.Fprintf(os.Stderr, "ln.Accept(): %s\n", err.Error())
				os.Exit(3)
			}
			rrs.AddStream(c)
		}
		fc := fnet.FrameConn(rrs)
		err = clockstation.Run(fc, time.Tick(50*time.Millisecond))
		if err != nil {
			fmt.Fprintf(os.Stderr, "clockstation.Run: %s\n", err.Error())
			os.Exit(2)
		}
	} else {
		rrs := fnet.NewStream(*frameSize)
		for i := 0; i < *numStreams; i++ {
			c, err := net.Dial("tcp", *addr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "net.Dial(%q): %s\n", *addr, err.Error())
				os.Exit(2)
			}
			rrs.AddStream(c)
		}
		fc := fnet.FrameConn(rrs)
		err := clockprinter.Run(fc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "clockprinter.Run: %s\n", err.Error())
			os.Exit(2)
		}
	}
}

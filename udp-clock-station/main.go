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

func main() {
	flag.Parse()
	if len(flag.Args()) != 0 || *listen == false && *addr == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *listen {
		localAddr, err := net.ResolveUDPAddr("udp", *addr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "net.ResolveUDPAddr(%q): %s\n", *addr, err.Error())
			os.Exit(2)
		}
		lc, err := net.ListenUDP("udp", localAddr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "net.Listen(%q): %s\n", *addr, err.Error())
			os.Exit(2)
		}
		_, remoteAddr, err := lc.ReadFromUDP(nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "lc.ReadFromUDP(): %s\n", err.Error())
			os.Exit(3)
		}
		lc.Close()
		c, err := net.DialUDP("udp", localAddr, remoteAddr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "net.DialUDP(\"udp\", %s, %s)", localAddr, remoteAddr, err.Error())
			os.Exit(4)
		}
		fc := fnet.Wrap(c, *frameSize)
		err = clockstation.Run(fc, time.Tick(50*time.Millisecond))
		if err != nil {
			fmt.Fprintf(os.Stderr, "clockstation.Run: %s\n", err.Error())
			os.Exit(5)
		}
	} else {
		remoteAddr, err := net.ResolveUDPAddr("udp", *addr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "net.ResolveUDPAddr(%q): %s\n", *addr, err.Error())
			os.Exit(2)
		}
		c, err := net.DialUDP("udp", nil, remoteAddr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "net.DialUDP(\"udp\", nil, %s)", remoteAddr, err.Error())
			os.Exit(3)
		}
		// TODO: resend dummy with exp. backoff until we get some response
		_, err = c.Write(nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "write UDP dummy to", remoteAddr, err.Error())
			os.Exit(4)
		}
		fc := fnet.Wrap(c, *frameSize)
		err = clockprinter.Run(fc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "clockprinter.Run: %s\n", err.Error())
			os.Exit(5)
		}
	}
}

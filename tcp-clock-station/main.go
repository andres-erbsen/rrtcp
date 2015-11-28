package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
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
		var cs *clockstation.ClockStation
		var fc fnet.FrameConn

		// Handle stop signals
		stop := make(chan os.Signal, 2)
		done := make(chan bool, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-stop
			cs.Stop()
			fc.Stop()
			done <- true
			fmt.Println("Stopped listener.")
			os.Exit(1)
		}()

		ln, err := net.Listen("tcp", *addr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "net.Listen(%q): %s\n", *addr, err.Error())
			os.Exit(2)
		}
		c, err := ln.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ln.Accept(): %s\n", err.Error())
			os.Exit(3)
		}
		fc = fnet.FromOrderedStream(c, *frameSize)
		cs = clockstation.NewStation(fc, time.Tick(50*time.Millisecond))
		err = cs.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "clockstation.Run: %s\n", err.Error())
			os.Exit(2)
		}
		// Wait for listener to be stopped
		<-done
	} else {
		var fc fnet.FrameConn

		// Handle stop signals
		stop := make(chan os.Signal, 2)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-stop
			fc.Stop()
			fmt.Println("Stopped dialer.")
			os.Exit(1)
		}()

		c, err := net.Dial("tcp", *addr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "net.Dial(%q): %s\n", *addr, err.Error())
			os.Exit(2)
		}
		fc = fnet.FromOrderedStream(c, *frameSize)
		err = clockprinter.Run(fc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "clockprinter.Run: %s\n", err.Error())
			os.Exit(2)
		}
	}
}

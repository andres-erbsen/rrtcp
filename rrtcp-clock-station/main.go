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
	var cs *clockstation.ClockStation

	// Handle stop signals
	// TODO: Is there a better place to put this?
	// TODO: Assumes cs, rrs have been defined
	// TODO: Is it better to os.Exit() or return?
	stop := make(chan os.Signal, 2)
	done := make(chan bool, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-stop
		cs.Stop()
		done <- true
		fmt.Println("Stopped listener.")
	}()

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "net.Listen(%q): %s\n", *addr, err.Error())
		return err
	}
	rrs := fnet.NewStream(*frameSize)

	for i := 0; i < *numStreams; i++ {
		c, err := ln.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ln.Accept(): %s\n", err.Error())
			return err
		}
		rrs.AddStream(c)
	}
	fc := fnet.FrameConn(rrs)
	defer fc.Stop()

	cs = clockstation.NewStation(fc, time.Tick(50*time.Millisecond))
	err = cs.Run(*frameSize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "clockstation.Run: %s\n", err.Error())
		return err
	}
	// Wait for listener to be stopped before returning
	<-done
	return nil
}

func dialer(frameSize *int, numStreams *int, addr *string) error {
	var fc fnet.FrameConn

	// Handle stop signals
	// TODO: Is there a better place to put this?
	// TODO: Assumes rrs has been defined
	stop := make(chan os.Signal, 2)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-stop
		// TODO: This should be a defer instead, when .Stop() can be called more than once freely
		fc.Stop()
		fmt.Println("Stopped dialer.")
	}()

	rrs := fnet.NewStream(*frameSize)
	for i := 0; i < *numStreams; i++ {
		c, err := net.Dial("tcp", *addr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "net.Dial(%q): %s\n", *addr, err.Error())
			return err
		}
		rrs.AddStream(c)
	}
	fc = fnet.FrameConn(rrs)

	err := clockprinter.Run(fc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "clockprinter.Run: %s\n", err.Error())
		return err
	}
	return nil
}

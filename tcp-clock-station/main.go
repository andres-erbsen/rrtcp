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
var duration = flag.Int("d", -1, "number of seconds to run program for, -1 means run forever")

func main() {
	flag.Parse()
	if len(flag.Args()) != 0 || *listen == false && *addr == "" {
		flag.Usage()
		os.Exit(1)
	}

	var stop chan os.Signal
	stop = make(chan os.Signal, 2)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if *duration != -1 {
			timer := time.NewTimer(time.Second * time.Duration(*duration))
			// Wait for the timer to end, then give the stop signal
			<-timer.C
			stop <- syscall.SIGINT
		}
	}()

	if *listen {
		err := listener(frameSize, addr, stop)
		if err != nil {
			os.Exit(2) // TODO: More appropriate per-case error number
		}
	} else {
		err := dialer(frameSize, addr, stop)
		if err != nil {
			os.Exit(2)
		}
	}
}

func listener(frameSize *int, addr *string, stop chan os.Signal) error {
	var cs *clockstation.ClockStation

	// Handle stop signals
	done := make(chan bool, 1)
	go func() {
		<-stop
		cs.Stop()
		done <- true
		fmt.Println("Stopped listener.")
	}()

	ln, err := net.Listen("tcp", *addr)
	defer ln.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "net.Listen(%q): %s\n", *addr, err.Error())
		return err
	}
	c, err := ln.Accept()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ln.Accept(): %s\n", err.Error())
		return err
	}
	fc := fnet.FromOrderedStream(c, *frameSize)
	defer fc.Stop()

	cs = clockstation.NewStation(fc, time.Tick(50*time.Millisecond))
	err = cs.Run(*frameSize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "clockstation.Run: %s\n", err.Error())
		return err
	}
	// Wait for listener to be stopped
	<-done
	return nil
}

func dialer(frameSize *int, addr *string, stop chan os.Signal) error {
	var fc fnet.FrameConn

	// Handle stop signals
	go func() {
		<-stop
		fc.Stop()
		fmt.Println("Stopped dialer.")
	}()

	c, err := net.Dial("tcp", *addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "net.Dial(%q): %s\n", *addr, err.Error())
		return err
	}
	fc = fnet.FromOrderedStream(c, *frameSize)
	err = clockprinter.Run(fc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "clockprinter.Run: %s\n", err.Error())
		fc.Stop()
		return err
	}
	return nil
}

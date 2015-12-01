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
		err := listener(*frameSize, *numStreams, *addr, stop)
		if err != nil {
			os.Exit(2)
		}
	} else {
		err := dialer(*frameSize, *numStreams, *addr, stop)
		if err != nil {
			os.Exit(2)
		}
	}
}

func listener(frameSize int, numStreams int, addr string, stop chan os.Signal) error {
	var cs *clockstation.ClockStation
	cs_running := false

	// Handle stop signals
	done := make(chan bool, 1)
	go func() {
		<-stop
		fmt.Println(cs)
		// TODO: Not sure if this boolean is the best solution
		if cs_running {
			cs.Stop()
		}
		done <- true
		fmt.Println("Stopped listener.")
	}()

	ln, err := net.Listen("tcp", addr)
	defer ln.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "net.Listen(%q): %s\n", addr, err.Error())
		return err
	}
	rr := fnet.NewRoundRobin(frameSize)

	for i := 0; i < numStreams; i++ {
		c, err := ln.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ln.Accept(): %s\n", err.Error())
			return err
		}
		fs := fnet.FromOrderedStream(c, frameSize)
		rr.AddConn(&fs)
	}
	fc := fnet.FrameConn(rr)
	defer fc.Stop()

	cs = clockstation.NewStation(fc, time.Tick(50*time.Millisecond))
	cs_running = true
	err = cs.Run(frameSize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "clockstation.Run: %s\n", err.Error())
		return err
	}
	// Wait for listener to be stopped before returning
	<-done
	return nil
}

func dialer(frameSize int, numStreams int, addr string, stop chan os.Signal) error {
	var fc fnet.FrameConn
	fc_started := false

	// Handle stop signals
	go func() {
		<-stop
		// TODO: This should be a defer instead, when .Stop() can be called more than once freely
		if fc_started {
			fc.Stop()
		}
		fmt.Println("Stopped dialer.")
	}()

	rr := fnet.NewRoundRobin(frameSize)
	for i := 0; i < numStreams; i++ {
		c, err := net.Dial("tcp", addr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "net.Dial(%q): %s\n", addr, err.Error())
			return err
		}
		fs := fnet.FromOrderedStream(c, frameSize)
		rr.AddConn(&fs)
	}
	fc = fnet.FrameConn(rr)
	fc_started = true

	err := clockprinter.Run(fc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "clockprinter.Run: %s\n", err.Error())
		return err
	}
	return nil
}

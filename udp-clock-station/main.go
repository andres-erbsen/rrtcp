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
var duration = flag.Int("d", 20, "number of seconds to run program for")

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
		timer := time.NewTimer(time.Second * time.Duration(*duration))
		// Wait for the timer to end, then give the stop signal
		<-timer.C
		stop <- syscall.SIGINT
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

	localAddr, err := net.ResolveUDPAddr("udp", *addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "net.ResolveUDPAddr(%q): %s\n", *addr, err.Error())
		return err
	}
	lc, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "net.Listen(%q): %s\n", *addr, err.Error())
		return err
	}
	_, remoteAddr, err := lc.ReadFromUDP(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "lc.ReadFromUDP(): %s\n", err.Error())
		return err
	}
	lc.Close()
	c, err := net.DialUDP("udp", localAddr, remoteAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "net.DialUDP(\"udp\", %s, %s)", localAddr, remoteAddr, err.Error())
		return err
	}
	fc := fnet.Wrap(c, *frameSize)
	defer fc.Stop()

	cs = clockstation.NewStation(fc, time.Tick(50*time.Millisecond))
	err = cs.Run(*frameSize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "clockstation.Run: %s\n", err.Error())
		return err
	}
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

	remoteAddr, err := net.ResolveUDPAddr("udp", *addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "net.ResolveUDPAddr(%q): %s\n", *addr, err.Error())
		return err
	}
	c, err := net.DialUDP("udp", nil, remoteAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "net.DialUDP(\"udp\", nil, %s)", remoteAddr, err.Error())
		return err
	}
	// TODO: resend dummy with exp. backoff until we get some response
	_, err = c.Write(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "write UDP dummy to", remoteAddr, err.Error())
		return err
	}
	fc = fnet.Wrap(c, *frameSize)
	err = clockprinter.Run(fc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "clockprinter.Run: %s\n", err.Error())
		fc.Stop()
		return err
	}
	return nil
}

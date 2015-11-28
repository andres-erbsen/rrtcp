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
		fc = fnet.Wrap(c, *frameSize)
		cs = clockstation.NewStation(fc, time.Tick(50*time.Millisecond))
		err = cs.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "clockstation.Run: %s\n", err.Error())
			os.Exit(5)
		}
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
		fc = fnet.Wrap(c, *frameSize)
		err = clockprinter.Run(fc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "clockprinter.Run: %s\n", err.Error())
			os.Exit(5)
		}
	}
}

package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/net/context"

	"github.com/andres-erbsen/rrtcp/clockprinter"
	"github.com/andres-erbsen/rrtcp/clockstation"
	"github.com/andres-erbsen/rrtcp/fnet"
)

var addr = flag.String("address", "", "address to connect to or listen at")
var listen = flag.Bool("l", false, "bind to the specified address and listen (default: connect)")
var duration = flag.Duration("d", 0, "duration to run program for")
var interval = flag.Duration("i", 50*time.Millisecond, "inter-packet interval")
var frameSize = flag.Int("s", 1024, "frame size")

func cancelOnSignal(ctx context.Context, sig ...os.Signal) context.Context {
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, sig...)

	ctx2, cancel := context.WithCancel(ctx)

	go func() {
		select {
		case <-ctx2.Done():
			return
		case <-signalCh:
			cancel()
		}
	}()

	return ctx2
}

func main() {
	fmt.Fprintln(os.Stderr, os.Args)
	flag.Parse()
	if len(flag.Args()) != 0 || *listen == false && *addr == "" {
		flag.Usage()
		os.Exit(1)
	}

	ctx := cancelOnSignal(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	if duration != nil && *duration != time.Duration(0) {
		ctx, _ = context.WithTimeout(ctx, *duration)
	}

	if *listen {
		if err := listener(ctx, *frameSize, *addr, *interval); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(2)
		}
	} else {
		if err := dialer(ctx, *frameSize, *addr); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(3)
		}
	}
	fmt.Fprintln(os.Stderr, "end main")
}

func listener(ctx context.Context, frameSize int, addr string, interval time.Duration) error {
	localAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return fmt.Errorf("net.ResolveUDPAddr(%q): %s\n", addr, err.Error())
	}
	lc, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		return fmt.Errorf("net.Listen(%q): %s\n", addr, err.Error())
	}
	go func() {
		<-ctx.Done()
		lc.Close()
	}()

	_, remoteAddr, err := lc.ReadFromUDP(nil)
	if err != nil {
		return fmt.Errorf("lc.ReadFromUDP(): %s\n", err.Error())
	}
	lc.Close()

	c, err := net.DialUDP("udp", localAddr, remoteAddr)
	if err != nil {
		return fmt.Errorf("net.DialUDP(\"udp\", nil, %s): %s", remoteAddr, err.Error())
	}

	if err := clockstation.Run(ctx, fnet.Wrap(c, frameSize), time.Tick(interval)); err != nil {
		return fmt.Errorf("clockstation.Run: %s\n", err.Error())
	}
	return nil
}

func dialer(ctx context.Context, frameSize int, addr string) error {
	remoteAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Errorf("net.ResolveUDPAddr(%q): %s\n", addr, err.Error())
	}
	c1, err := net.DialUDP("udp", nil, remoteAddr)
	if err != nil {
		return fmt.Errorf("net.DialUDP(\"udp\", nil, %s): %s", remoteAddr, err.Error())
	}
	go func() {
		<-ctx.Done()
		c1.Close()
	}()

	c2, err := net.DialUDP("udp", nil, remoteAddr)
	// TODO: resend dummy with exp. backoff until we get some response
	if _, err = c2.Write(nil); err != nil {
		return fmt.Errorf("write UDP dummy to %s: %s", remoteAddr, err.Error())
	}
	go func() {
		<-ctx.Done()
		c2.Close()
	}()

	if err := clockprinter.Run(ctx, fnet.Wrap(c2, frameSize)); err != nil {
		fmt.Errorf("clockprinter.Run: %s\n", err.Error())
	}
	return nil
}

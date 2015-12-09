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
var frameSize = flag.Int("s", 1024, "frame size")
var numStreams = flag.Int("n", 5, "number of streams")
var duration = flag.Duration("d", 0, "duration to run program for")

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
		err := listener(ctx, *frameSize, *numStreams, *addr)
		if err != nil {
			os.Exit(2)
		}
	} else {
		err := dialer(ctx, *frameSize, *numStreams, *addr)
		if err != nil {
			os.Exit(2)
		}
	}
}

func listener(ctx context.Context, frameSize int, numStreams int, addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "net.Listen(%q): %s\n", addr, err.Error())
		return err
	}
	go func() {
		<-ctx.Done()
		ln.Close()
	}()

	rr := fnet.NewRoundRobin(frameSize)
	go func() {
		<-ctx.Done()
		rr.Close()
	}()

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

	if err = clockstation.Run(ctx, fc, time.Tick(50*time.Millisecond)); err != nil {
		return fmt.Errorf("clockstation.Run: %s\n", err.Error())
	}
	return nil
}

func dialer(ctx context.Context, frameSize int, numStreams int, addr string) error {
	rr := fnet.NewRoundRobin(frameSize)
	go func() {
		<-ctx.Done()
		rr.Close()
	}()

	for i := 0; i < numStreams; i++ {
		c, err := net.Dial("tcp", addr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "net.Dial(%q): %s\n", addr, err.Error())
			return err
		}
		fs := fnet.FromOrderedStream(c, frameSize)
		rr.AddConn(&fs)
	}
	fc := fnet.FrameConn(rr)

	if err := clockprinter.Run(ctx, fc); err != nil {
		fmt.Errorf("clockprinter.Run: %s\n", err.Error())
	}
	return nil
}

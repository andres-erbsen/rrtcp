package fnet

import (
	"fmt"
	"net"
)

type netConn struct {
	c         net.Conn
	frameSize int
}

// Wrap wraps a net.Conn that already implements framing
func Wrap(c net.Conn, frameSize int) FrameConn {
	return netConn{c, frameSize}
}

func (c netConn) FrameSize() int {
	return c.frameSize
}

func (c netConn) SendFrame(bs []byte) error {
	if len(bs) != c.frameSize {
		return fmt.Errorf("Tried to send byte of length %d instead of %d", len(bs), c.frameSize)
	}
	_, err := c.c.Write(bs)
	return err
}

func (c netConn) RecvFrame(bs []byte) error {
	n, err := c.c.Read(bs)
	if err != nil {
		return err
	}
	if n != len(bs) {
		return fmt.Errorf("received frame of incorrect length")
	}
	return nil
}

func (c netConn) Close() error {
	return c.c.Close()
}

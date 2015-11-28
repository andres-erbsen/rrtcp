package fnet

import "net"

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
	_, err := c.c.Write(bs)
	return err
}

func (c netConn) RecvFrame(bs []byte) (int, error) {
	return c.c.Read(bs)
}

func (c netConn) Stop() {
	c.c.Close()
}

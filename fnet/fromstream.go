package fnet

import (
	"encoding/binary"
	"io"
	"net"
)

type framedStream struct {
	c         net.Conn
	frameSize int
}

// FromOrderedStream wraps a net.Conn in a framing layer.
// PRE: c :-> net.Conn
// RET: ret :-> FrameConn{FrameSize = frameSize}
func FromOrderedStream(c net.Conn, frameSize int) FrameConn {
	return &framedStream{c, frameSize}
}

// FrameSize implements FrameConn.FrameSize
func (fs *framedStream) FrameSize() int {
	return fs.frameSize
}

// SendFrame implements FrameConn.SendFrame
func (fs *framedStream) SendFrame(b []byte) error {
	b2 := make([]byte, binary.MaxVarintLen64+len(b))
	i := binary.PutUvarint(b2, uint64(len(b)))
	if copy(b2[i:], b) != len(b) {
		panic("length accounting failed")
	}
	_, err := fs.c.Write(b2)
	return err
}

type byteReader struct{ io.Reader }

func (r byteReader) ReadByte() (byte, error) {
	var ret [1]byte
	_, err := io.ReadFull(r, ret[:])
	return ret[0], err
}

// RecvFrame implements FrameConn.RecvFrame
func (fs *framedStream) RecvFrame(b []byte) (int, error) {
	sz, err := binary.ReadUvarint(byteReader{fs.c})
	if err != nil {
		return 0, err
	}
	return io.ReadFull(fs.c, b[:sz])
}

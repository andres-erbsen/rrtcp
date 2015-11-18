package fnet

import (
	"net"
)

const numConn = 20
const recBufSize = 1 * numConn // Received packet buffer is a linear function of the number of connections

type rrStream struct {
	pool       [numConn]*framedStream
	frameSize  int
	nextStream int         // The index of the next stream to send a packet on
	rec        chan []byte // Queue of received packets
}

// Makes a round robin Stream given a pool of TCP connections and a frameSize
// PRE: connPool :-> [numConn]net.Conn
// RET: ret :-> FrameConn{FrameSize = frameSize}
func FromConnPool(connPool [numConn]net.Conn, frameSize int) FrameConn {
	var streamPool [numConn]*framedStream
	for i := 0; i < numConn; i++ {
		streamPool[i] = &framedStream{connPool[i], frameSize}
	}
	rrs := &rrStream{streamPool, frameSize, 0, make(chan []byte, recBufSize)}
	// Start a new thread for listening to every connection
	for i := 0; i < numConn; i++ {
		// TODO: Do something with possible errors
		// TODO: WaitGroups so we can stop the stream
		go rrs.Listen(streamPool[i])
	}
	return rrs
}

// FrameSize implements FrameConn.FrameSize
func (rrs *rrStream) FrameSize() int {
	return rrs.frameSize
}

// Listen for incoming packets and add them to the received queue
func (rrs *rrStream) Listen(fs *framedStream) error {
	// TODO: Safe stop
	for {
		buf := make([]byte, rrs.frameSize)
		sz, err := fs.RecvFrame(buf)
		if err != nil {
			return err
		}
		rrs.rec <- buf[:sz]
	}
}

// SendFrame implements FrameConn.SendFrame
func (rrs *rrStream) SendFrame(b []byte) error {
	fs := rrs.pool[rrs.nextStream]
	err := fs.SendFrame(b)
	// TODO: Should we actually move up to the next stream if there's an error?
	rrs.nextStream = (rrs.nextStream + 1) % numConn // Get the next round-robin index
	return err
}

func (rrs *rrStream) RecvFrame(b []byte) (int, error) {
	frame := <-rrs.rec
	copy(b[:len(frame)], frame)

	return len(frame), nil
}

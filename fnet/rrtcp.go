package fnet

import (
	"net"
	"sync"
)

const recBufSize = 20 // This is a total guess as to a reasonable buffer size for our receive channel

type rrStream struct {
	pool       []*framedStream
	poolLock   sync.Mutex // Lock for changing the pool and pool related values (numStreams)
	frameSize  int
	numStreams int
	nextStream int         // The index of the next stream to send a packet on
	rec        chan []byte // Queue of received packets
	wg         sync.WaitGroup
	stop       bool
}

func (rrs *rrStream) AddStream(conn net.Conn) {
	stream := &framedStream{conn, rrs.frameSize}

	rrs.poolLock.Lock()
	rrs.numStreams++
	rrs.pool = append(rrs.pool, stream)
	rrs.poolLock.Unlock()
	// Start a new thread for listening to every connection
	rrs.wg.Add(1)
	go rrs.Listen(stream, rrs.numStreams-1)
}

func NewStream(frameSize int) *rrStream {
	var streamPool []*framedStream
	var wg sync.WaitGroup
	var lock sync.Mutex
	rrs := &rrStream{streamPool, lock, frameSize, 0, 0, make(chan []byte, recBufSize), wg, false}
	return rrs
}

// FrameSize implements FrameConn.FrameSize
func (rrs *rrStream) FrameSize() int {
	return rrs.frameSize
}

func (rrs *rrStream) Stop() {
	rrs.stop = true
	for _, stream := range rrs.pool {
		stream.c.Close()
	}
	rrs.wg.Wait()
}

// Listen for incoming packets and add them to the received queue
func (rrs *rrStream) Listen(fs *framedStream, index int) {
	defer rrs.wg.Done()
	for {
		buf := make([]byte, rrs.frameSize)
		sz, err := fs.RecvFrame(buf)
		if err != nil {
			if rrs.stop { // Stop this thread
				return
			} else {
				// Remove the stream if the connection is sad
				rrs.RemoveStream(fs, index)
				return
			}
		}
		rrs.rec <- buf[:sz]
	}
}

func (rrs *rrStream) RemoveStream(fs *framedStream, index int) {
	fs.c.Close()
	rrs.poolLock.Lock()
	rrs.numStreams--
	rrs.pool = append(rrs.pool[:index], rrs.pool[index+1:]...)
	rrs.poolLock.Unlock()
}

// SendFrame implements FrameConn.SendFrame
func (rrs *rrStream) SendFrame(b []byte) error {
	fs := rrs.pool[rrs.nextStream]
	err := fs.SendFrame(b)
	// TODO: Should we actually move up to the next stream if there's an error?
	rrs.nextStream = (rrs.nextStream + 1) % rrs.numStreams // Get the next round-robin index
	return err
}

// RecvFrame implements FrameConn.RecvFrame
// It pulls the next frame out of the rec channel
// This method should be running continously to prevent blocking on the rec chan
func (rrs *rrStream) RecvFrame(b []byte) (int, error) {
	frame := <-rrs.rec
	copy(b[:len(frame)], frame)

	return len(frame), nil
}

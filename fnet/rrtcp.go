package fnet

import (
	"errors"
	"net"
	"sync"
)

const recBufSize = 20 // This is a total guess as to a reasonable buffer size for our receive channel

type RrStream struct {
	pool       []*framedStream
	poolLock   sync.Mutex // Lock for changing the pool and pool related values (numStreams)
	frameSize  int
	numStreams int
	nextStream int         // The index of the next stream to send a packet on
	rec        chan []byte // Queue of received packets
	wg         sync.WaitGroup
	stopCh     chan struct{}
}

func (rrs *RrStream) AddStream(conn net.Conn) {
	stream := &framedStream{conn, rrs.frameSize}

	rrs.poolLock.Lock()
	rrs.numStreams++
	rrs.pool = append(rrs.pool, stream)
	rrs.poolLock.Unlock()
	// Start a new thread for listening to every connection
	rrs.wg.Add(1)
	go rrs.listen(stream)
}

func NewStream(frameSize int) *RrStream {
	var streamPool []*framedStream
	var wg sync.WaitGroup
	var lock sync.Mutex
	rrs := &RrStream{streamPool, lock, frameSize, 0, 0, make(chan []byte, recBufSize), wg, make(chan struct{})}
	return rrs
}

// FrameSize implements FrameConn.FrameSize
func (rrs *RrStream) FrameSize() int {
	return rrs.frameSize
}

// TODO: Make it so that this can be called more than once freely
func (rrs *RrStream) Stop() {
	close(rrs.stopCh)
	for _, stream := range rrs.pool {
		stream.c.Close()
	}
	rrs.wg.Wait()
}

// listen for incoming packets and add them to the received queue
func (rrs *RrStream) listen(fs *framedStream) {
	defer rrs.wg.Done()
	for {
		buf := make([]byte, rrs.frameSize)
		sz, err := fs.RecvFrame(buf)
		if err != nil {
			select {
			case <-rrs.stopCh: // Stop this thread
				return
			default:
				// Remove the stream if the connection is sad
				rrs.RemoveStream(fs)
				return
			}
		}
		rrs.rec <- buf[:sz]
	}
}

// TODO: Implement this more efficiently
func (rrs *RrStream) RemoveStream(fs *framedStream) {
	fs.c.Close()
	rrs.poolLock.Lock()

	// Get index of stream
	// TODO: Exception if doesn't exist at all
	var fsIndex int
	for index, stream := range rrs.pool {
		if stream == fs {
			fsIndex = index
			break
		}
	}

	rrs.numStreams--
	if rrs.nextStream >= rrs.numStreams {
		rrs.nextStream = 0
	}
	rrs.pool = append(rrs.pool[:fsIndex], rrs.pool[fsIndex+1:]...)
	rrs.poolLock.Unlock()
}

// SendFrae implements FrameConn.SendFrame
func (rrs *RrStream) SendFrame(b []byte) error {
	rrs.poolLock.Lock()
	if rrs.numStreams == 0 {
		return errors.New("No streams to send packets on.")
	}
	fs := rrs.pool[rrs.nextStream]
	err := fs.SendFrame(b)
	// TODO: Should we actually move up to the next stream if there's an error?
	rrs.nextStream = (rrs.nextStream + 1) % rrs.numStreams // Get the next round-robin index
	rrs.poolLock.Unlock()
	return err
}

// RecvFrame implements FrameConn.RecvFrame
// It pulls the next frame out of the rec channel
// This method should be running continously to prevent blocking on the rec chan
func (rrs *RrStream) RecvFrame(b []byte) (int, error) {
	for {
		select {
		case <-rrs.stopCh: // Stop this thread
			return 0, errors.New("Stream stopped.")
		case frame := <-rrs.rec:
			copy(b[:len(frame)], frame)

			return len(frame), nil
		}
	}
}

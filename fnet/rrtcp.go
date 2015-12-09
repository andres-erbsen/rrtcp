package fnet

import (
	"errors"
	"fmt"
	"sync"
)

const recvBufSize = 20 // This is a total guess as to a reasonable buffer size for our recveive channel

type RoundRobin struct {
	pool      []FrameConn
	poolLock  sync.Mutex // Lock for changing the pool and pool related values (numConn)
	frameSize int
	numConn   int
	nextConn  int         // The index of the next stream to send a packet on
	recv      chan []byte // Queue of received packets
	wg        sync.WaitGroup
	stopCh    chan struct{}
}

var _ FrameConn = (*RoundRobin)(nil)

func (rr *RoundRobin) AddConn(fc FrameConn) {
	rr.poolLock.Lock()
	defer rr.poolLock.Unlock()
	rr.numConn++
	rr.pool = append(rr.pool, fc)

	// Start a new thread for listening to every connection
	rr.wg.Add(1)
	go rr.listen(fc)
}

func NewRoundRobin(frameSize int) *RoundRobin {
	var conn []FrameConn
	var wg sync.WaitGroup
	var lock sync.Mutex
	rr := &RoundRobin{
		pool:      conn,
		poolLock:  lock,
		frameSize: frameSize,
		numConn:   0,
		nextConn:  0,
		recv:      make(chan []byte, recvBufSize),
		wg:        wg,
		stopCh:    make(chan struct{}),
	}
	return rr
}

// FrameSize implements FrameConn.FrameSize
func (rr *RoundRobin) FrameSize() int {
	return rr.frameSize
}

// TODO: Make it so that this can be called more than once freely
func (rr *RoundRobin) Close() error {
	close(rr.stopCh)
	var err error
	for _, conn := range rr.pool {
		err = conn.Close()
	}
	rr.wg.Wait()
	return err
}

// listen for incoming packets and add them to the received queue
func (rr *RoundRobin) listen(fc FrameConn) {
	defer rr.wg.Done()
	for {
		buf := make([]byte, rr.frameSize)
		if err := fc.RecvFrame(buf); err != nil {
			select {
			case <-rr.stopCh: // Stop this thread
				fmt.Printf("RoundRobin.listen: %v", err)
				return
			default:
				// Remove the stream if the connection is sad
				rr.RemoveConn(fc)
				return
			}
		}
		rr.recv <- buf[:]
	}
}

// TODO: Implement this more efficiently
func (rr *RoundRobin) RemoveConn(fc FrameConn) {
	fc.Close()
	rr.poolLock.Lock()

	// Get index of stream
	// TODO: Exception if doesn't exist at all
	var fcIndex int
	for index, conn := range rr.pool {
		if conn == fc {
			fcIndex = index
			break
		}
	}

	rr.numConn--
	if rr.nextConn >= rr.numConn {
		rr.nextConn = 0
	}
	rr.pool = append(rr.pool[:fcIndex], rr.pool[fcIndex+1:]...)
	rr.poolLock.Unlock()
}

// SendFrame implements FrameConn.SendFrame
func (rr *RoundRobin) SendFrame(b []byte) error {
	rr.poolLock.Lock()
	if rr.numConn == 0 {
		return errors.New("No streams to send packets on.")
	}
	fc := rr.pool[rr.nextConn]
	err := fc.SendFrame(b)
	if err != nil {
		rr.poolLock.Unlock()
		rr.RemoveConn(fc)
		return err
	} else {
		rr.nextConn = (rr.nextConn + 1) % rr.numConn // Get the next round-robin index
		rr.poolLock.Unlock()
		return nil
	}
}

// RecvFrame implements FrameConn.RecvFrame
// It pulls the next frame out of the recv channel
// This method should be running continously to prevent blocking on the recv chan
func (rr *RoundRobin) RecvFrame(b []byte) error {
	for {
		select {
		case <-rr.stopCh: // Stop this thread
			return errors.New("Stream stopped.")
		case frame := <-rr.recv:
			copy(b, frame)
			return nil
		}
	}
}

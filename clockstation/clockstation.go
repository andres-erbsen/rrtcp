package clockstation

import (
	"encoding/binary"
	"time"

	"github.com/andres-erbsen/rrtcp/fnet"
)

type clockStation struct {
	fc      fnet.FrameConn
	tick    <-chan time.Time
	stopCh  chan struct{}
	stopped chan struct{}
}

// Start starts a new clock station. The connection argument is wrapped.
// PRE: fc :-> fnet.FrameConn
// EFF: sends each value from ticker over fc (nanoseconds in little-endian 64-bit signed unix epoch format) until stop is called
func Run(fc fnet.FrameConn, tick <-chan time.Time) error {
	cs := clockStation{fc, tick, make(chan struct{}), make(chan struct{})}
	return cs.run()
}

func (cs *clockStation) run() error {
	defer close(cs.stopped)
	var b [8]byte
	for {
		select {
		case t := <-cs.tick:
			ns := t.UnixNano()
			binary.LittleEndian.PutUint64(b[:], uint64(ns))
			if err := cs.fc.SendFrame(b[:]); err != nil {
				return err
			}
		case <-cs.stopCh:
			return nil
		}
	}
}

func (cs *clockStation) Stop() {
	close(cs.stopCh)
	<-cs.stopped
}

package clockstation

import (
	"encoding/binary"
	"fmt"
	"golang.org/x/net/context"
	"time"

	"github.com/andres-erbsen/rrtcp/fnet"
)

// Start starts a new clock station. The connection argument is wrapped.
// PRE: fc :-> fnet.FrameConn
// POST: ret = error OR nil
// EFF: sends each value from ticker over fc (nanoseconds in little-endian 64-bit signed unix epoch format) until stop is called
func Run(ctx context.Context, fc fnet.FrameConn, tick <-chan time.Time) error {
	var b = make([]byte, fc.FrameSize())
	for {
		select {
		case t := <-tick:
			ns := t.UnixNano()
			binary.LittleEndian.PutUint64(b[:8], uint64(ns))
			if err := fc.SendFrame(b[:]); err != nil {
				return err
			}
			fmt.Printf("%d\n", uint64(ns))
			//		fmt.Printf("Sent frame %v", ns) yo
		case <-ctx.Done():
			return nil
		}
	}
}

package clockprinter

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/andres-erbsen/rrtcp/fnet"
)

func Run(fc fnet.FrameConn) error {
	bs := make([]byte, fc.FrameSize())
	for {
		n, err := fc.RecvFrame(bs)
		if err != nil {
			return err
		}
		if n < fc.FrameSize() {
			return fmt.Errorf("frame too small (got %d, wanted %d", n, fc.FrameSize())
		}
		fmt.Printf("%d %d\n", int64(binary.LittleEndian.Uint64(bs[:8])), time.Now().UnixNano())
	}
}

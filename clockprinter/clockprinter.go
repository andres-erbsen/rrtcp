package clockprinter

import (
	"encoding/binary"
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/andres-erbsen/rrtcp/fnet"
)

func Run(ctx context.Context, fc fnet.FrameConn) error {
	bs := make([]byte, fc.FrameSize())
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if err := fc.RecvFrame(bs); err != nil {
			return err
		}
		fmt.Printf("%d %d\n", int64(binary.LittleEndian.Uint64(bs[:8])), time.Now().UnixNano())
	}
}

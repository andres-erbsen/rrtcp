package fnet

type FrameConn interface {
	FrameSize() int

	// SendFrame sends a bounded-size frame over the connection.
	// PRE: b :->[] bs, len(bs) = FrameSize
	// RET: b :->[] xs, len(xs) = FrameSize
	// EFF: if err = nil then SendFrame(bs) else (SendFrame(xs) OR NoEffects)
	SendFrame(b []byte) error // only frames of valid size or less

	// RecvFrame receives a bounded-size fram over the connection
	// PRE: b :->[] mm, len(mm) = FrameSize.
	// RET: b :->[] bs, len(bs) = FrameSize
	// EFF: if err = nil then RecvFrame(firstn n received) else (RecvFrame(firstn n bs) OR NoEffects)
	RecvFrame(b []byte) (n int, err error)
}

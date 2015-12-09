package udp4rrtcp

import (
	"fmt"
	"net"

	"github.com/andres-erbsen/rrtcp/fnet"
	"golang.org/x/net/context"
)

func Start(ctx context.Context, udpSourceListen, udpDestinationSend string, carrier fnet.FrameConn) (*UDP4RRTCP, error) {
	listenAddr, err := net.ResolveUDPAddr("udp", udpSourceListen)
	if err != nil {
		return nil, err
	}
	sendAddr, err := net.ResolveUDPAddr("udp", udpDestinationSend)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		return nil, err
	}

	ret := &UDP4RRTCP{conn, sendAddr, carrier}
	go ret.udp4rrtcpSendLoop()
	go ret.udp4rrtcpRecvLoop()
	return ret, nil
}

func (u *UDP4RRTCP) Stop() error {
	if err := u.carrier.Close(); err != nil {
		return err
	}
	return u.conn.Close()
}

type UDP4RRTCP struct {
	conn     *net.UDPConn
	sendAddr *net.UDPAddr
	carrier  fnet.FrameConn
}

func (u *UDP4RRTCP) udp4rrtcpSendLoop() {
	b := make([]byte, u.carrier.FrameSize())
	for {
		n, _, err := u.conn.ReadFrom(b[:])
		if err != nil {
			fmt.Printf("%s\n", err)
			continue
		}
		if n != len(b) {
			fmt.Printf("udp->rrtcp: expected %d bytes, got %d\n", len(b), n)
		}
		if err := u.carrier.SendFrame(b); err != nil {
			fmt.Printf("%s\n", err)
		}
	}
}

func (u *UDP4RRTCP) udp4rrtcpRecvLoop() {
	b := make([]byte, u.carrier.FrameSize())
	for {
		if err := u.carrier.RecvFrame(b); err != nil {
			fmt.Printf("%s\n", err)
		}
		if n, err := u.conn.WriteTo(b, u.sendAddr); err != nil || n != len(b) {
			fmt.Printf("udpconn.WriteTo returned %v, %v\n", n, err)
		}

	}
}

package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"os"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/proxy"

	"github.com/andres-erbsen/rrtcp/clockprinter"
	"github.com/andres-erbsen/rrtcp/clockstation"
	"github.com/andres-erbsen/rrtcp/fnet"
	"github.com/andres-erbsen/torch"
	"github.com/andres-erbsen/torch/directory"
	"github.com/andres-erbsen/torch/nd"
)

var identifiable = flag.Bool("i", false, "skip TOR anonymization")
var numStreams = flag.Int("n", 1, "number of streams")
var udpSourceListen = flag.String("source", ":", "UDP address to listen on and read packets from")
var udpDestinationSend = flag.String("sink", "", "UDP to send the received packets to")

type hex32byte [32]byte

func (h *hex32byte) String() string     { return fmt.Sprintf("%x", *h) }
func (h *hex32byte) Set(s string) error { _, err := fmt.Sscanf("%x", s, h[:]); return err }

func main() {
	var id [32]byte
	flag.Var((*hex32byte)(&id), "rend", "hex ID of the OR through which the connection is created. If zero, computed from seed.")
	flag.Parse()

	if flag.NArg() != 1 || len(flag.Arg(0)) == 0 {
		os.Stderr.Write([]byte("USAGE: please specify a shared seed as the first argument\n"))
		flag.Usage()
		os.Exit(3)
	}
	seed := flag.Arg(0)

	if id == [32]byte{} {
		id = sha256.Sum256(append([]byte("TF_EXPAND_SEED_R"), []byte(seed)...))
	}
	ctx := context.Background()
	err := run(ctx, &id, []byte(seed), *identifiable, *numStreams, *udpSourceListen, *udpDestinationSend)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, rendID *[32]byte, seed []byte, identifiable bool, numStreams int, udpSourceListen, udpDestinationSend string) error {
	tr, err := torch.New(ctx, proxy.FromEnvironment())
	if err != nil {
		return err
	}
	defer tr.Stop()

	rendNode := tr.WithDirectory(func(dir *directory.Directory) interface{} {
		shortInfos := make([]*directory.ShortNodeInfo, 0, len(dir.Routers))
		for _, n := range dir.Routers {
			if n.Fast && n.Running && n.Stable && n.Valid && (n.Port == 443 || n.Port == 80) {
				shortInfos = append(shortInfos, n.ShortNodeInfo)
			}
		}
		rendShort := nd.Pick(shortInfos, rendID)
		for _, r := range dir.Routers {
			if r.ShortNodeInfo == rendShort {
				return r
			}
		}
		panic("unreachable")
	}).(*directory.NodeInfo)

	// when we talk to a website over TOR, the website knows what the third hop
	// is. Therefore our converstion partner may as well.
	nodes := make([]*directory.NodeInfo, 0, 3)
	if !identifiable {
		nodes = append(nodes, tr.Pick(weigh))
		nodes = append(nodes, tr.Pick(weigh))
	}
	nodes = append(nodes, rendNode)

	rr := fnet.NewRoundRobin(nd.FRAMESIZE)
	defer rr.Close()

	sending := false

	for i := 0; i < numStreams; i++ {
		seed_i := sha256.Sum256(append(append([]byte{}, seed...), byte(i)))
		tc1, c1, err := torch.BuildCircuit(ctx, proxy.FromEnvironment(), nodes)
		if err != nil {
			return err
		}
		defer tc1.Close()
		tc2, c2, err := torch.BuildCircuit(ctx, proxy.FromEnvironment(), nodes)
		if err != nil {
			return err
		}
		defer tc2.Close()

		ndc, err := nd.Handshake(ctx, c1, c2, seed_i[:])
		// Arbitrarily choose one of the sides of the TOR connection to be the sender of the test data
		// By choosing the first ndc.bit
		if i == 0 {
			sending = ndc.Bit
		}
		if err != nil {
			return err
		}
		rr.AddConn(ndc)
	}

	println("connected and running")

	if sending {
		if err = clockstation.Run(ctx, rr, time.Tick(50*time.Millisecond)); err != nil {
			return fmt.Errorf("clockstation.Run: %s\n", err.Error())
		}
	} else {
		if err = clockprinter.Run(ctx, rr); err != nil {
			fmt.Errorf("clockprinter.Run: %s\n", err.Error())
		}
	}
	return nil

	<-ctx.Done()
	return nil
}

func weigh(w *directory.BandwidthWeights, n *directory.NodeInfo) int64 {
	if n.Fast && n.Running && n.Stable && n.Valid && (n.Port == 443 || n.Port == 80) {
		return w.ForRelay.Weigh(n)
	}
	return 0
}

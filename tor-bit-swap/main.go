package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/net/proxy"

	"github.com/andres-erbsen/rrtcp/fnet"
	"github.com/andres-erbsen/torch"
	"github.com/andres-erbsen/torch/directory"
	"github.com/andres-erbsen/torch/nd"
)

var identifiable = flag.Bool("i", false, "skip TOR anonymization")

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
	run(ctx, &id, []byte(seed), *identifiable)
}

func run(ctx context.Context, rendID *[32]byte, seed []byte, identifiable bool) error {
	tr, err := torch.New(ctx, proxy.FromEnvironment())
	if err != nil {
		return err
	}
	defer tr.Stop()

	rendNode := tr.WithDirectory(func(dir *directory.Directory) interface{} {
		shortInfos := make([]*directory.ShortNodeInfo, len(dir.Routers))
		for i, r := range dir.Routers {
			shortInfos[i] = r.ShortNodeInfo
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
		nodes = append(nodes, tr.Pick(weighRelayWith, nil))
		nodes = append(nodes, tr.Pick(weighRelayWith, nil))
	}
	nodes = append(nodes, rendNode)

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

	ndc, err := nd.Handshake(ctx, c1, c2, seed)
	var _ fnet.FrameConn = (ndc)

	fmt.Printf("%#v\n", ndc.Bit)

	return nil
}

func weighRelayWith(w *directory.BandwidthWeights, n *directory.NodeInfo) int64 {
	return w.ForRelay.Weigh(n)
}

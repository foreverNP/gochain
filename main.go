package main

import (
	"context"
	"time"

	"github.com/foreverNP/gochain/node"
	"github.com/foreverNP/gochain/proto"

	"google.golang.org/grpc"
)

func main() {
	makeNode(":3000", []string{})
	time.Sleep(1 * time.Second)
	makeNode(":4000", []string{":3000"})

	//go func() {
	//	for {
	//		time.Sleep(2 * time.Second)
	//		makeTransaction()
	//	}
	//}()

	select {}
}

func makeNode(lnAddr string, bootstrapAddrs []string) *node.Node {
	n := node.New()
	go n.Start(lnAddr)
	if len(bootstrapAddrs) > 0 {
		if err := n.BootstrapNetwork(bootstrapAddrs); err != nil {
			panic(err)
		}
	}

	return n
}

func makeTransaction() {
	client, err := grpc.Dial(":3000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	c := proto.NewNodeClient(client)

	ver := &proto.Version{
		Version:    "0.0.1",
		Height:     100,
		ListenAddr: ":3000",
	}

	_, err = c.Handshake(context.TODO(), ver)
	if err != nil {
		panic(err)
	}
}

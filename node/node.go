package node

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/foreverNP/gochain/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type Node struct {
	version string
	lnAddr  string
	logger  *zap.SugaredLogger

	peerLock sync.RWMutex
	peers    map[proto.NodeClient]*proto.Version

	proto.UnimplementedNodeServer
}

func New() *Node {
	loggerConf := zap.NewDevelopmentConfig()
	loggerConf.EncoderConfig.TimeKey = ""
	logger, _ := loggerConf.Build()
	return &Node{
		peers:   make(map[proto.NodeClient]*proto.Version),
		version: "0.0.1",
		logger:  logger.Sugar(),
	}
}

func (n *Node) Start(lnAddr string) error {
	n.lnAddr = lnAddr
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	ln, err := net.Listen("tcp", lnAddr)
	if err != nil {
		return err
	}
	proto.RegisterNodeServer(grpcServer, n)

	n.logger.Infow("node started...", "addr", lnAddr)

	return grpcServer.Serve(ln)
}

func (n *Node) addPeer(c proto.NodeClient, v *proto.Version) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()

	n.logger.Infof("[%s]: new peer connected [%s]", n.lnAddr, v.GetListenAddr())

	n.peers[c] = v
}

func (n *Node) removePeer(c proto.NodeClient) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()
	delete(n.peers, c)
}

func (n *Node) BootstrapNetwork(addrs []string) error {
	for _, addr := range addrs {
		c, err := makeNodeClient(addr)
		if err != nil {
			return err
		}

		v, err := c.Handshake(context.Background(), n.getVersion())
		if err != nil {
			n.logger.Errorf("failed to connect to peer: %s", err)
			continue
		}

		n.addPeer(c, v)
	}

	return nil
}

func (n *Node) getVersion() *proto.Version {
	return &proto.Version{
		Version:    n.version,
		Height:     0,
		ListenAddr: n.lnAddr,
	}
}

func (n *Node) Handshake(ctx context.Context, v *proto.Version) (*proto.Version, error) {
	c, err := makeNodeClient(v.GetListenAddr())
	if err != nil {
		return nil, err
	}

	n.addPeer(c, v)

	return n.getVersion(), nil
}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*proto.Ack, error) {
	p, _ := peer.FromContext(ctx)
	fmt.Println("received transaction:", p)
	return &proto.Ack{}, nil
}

func makeNodeClient(lnAddr string) (proto.NodeClient, error) {
	conn, err := grpc.Dial(lnAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return proto.NewNodeClient(conn), nil
}

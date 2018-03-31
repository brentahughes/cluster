package node

import (
	"context"
	"log"
	"net"

	"github.com/bah2830/cluster/service"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func (n *Node) Checkin(ctx context.Context, in *service.Empty) (*service.NodeDetails, error) {
	return n.GetNodeDetails(), nil
}

func (n *Node) Execute(ctx context.Context, in *service.ExecutionRequest) (*service.ExecutionResponse, error) {
	return &service.ExecutionResponse{}, nil
}

func (n *Node) Details(ctx context.Context, in *service.Empty) (*service.NodeDetailsVerbose, error) {
	return &service.NodeDetailsVerbose{}, nil
}

func (n *Node) startAPIServer() {
	log.Printf("Starting node service on port %s", viper.GetString("rpc.port"))

	lis, err := net.Listen("tcp", ":"+viper.GetString("rpc.port"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Creates a new gRPC server
	s := grpc.NewServer()
	service.RegisterNodeServer(s, n)
	s.Serve(lis)
}

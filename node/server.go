package node

import (
	"context"
	"log"
	"net"

	"github.com/bah2830/cluster/service"
	"google.golang.org/grpc"
)

func (n *node) Checkin(ctx context.Context, in *service.Empty) (*service.NodeDetails, error) {
	return &service.NodeDetails{}, nil
}

func (n *node) Execute(ctx context.Context, in *service.ExecutionRequest) (*service.ExecutionResponse, error) {
	return &service.ExecutionResponse{}, nil
}

func startServer(port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Creates a new gRPC server
	s := grpc.NewServer()
	service.RegisterNodeServer(s, &node{})
	s.Serve(lis)
}

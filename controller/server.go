package controller

import (
	"context"
	"log"
	"net"

	"github.com/davecgh/go-spew/spew"

	"github.com/bah2830/cluster/service"
	"google.golang.org/grpc"
)

func (c *controller) Checkin(ctx context.Context, in *service.NodeDetails) (*service.GenericReply, error) {
	spew.Dump(in)
	return &service.GenericReply{}, nil
}

func (c *controller) ExecResponse(ctx context.Context, in *service.ExecutionResponse) (*service.Empty, error) {
	return &service.Empty{}, nil
}

func startServer(port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Creates a new gRPC server
	s := grpc.NewServer()
	service.RegisterControllerServer(s, &controller{})
	s.Serve(lis)
}

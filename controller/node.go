package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/bah2830/cluster/service"
	"github.com/boltdb/bolt"
	"google.golang.org/grpc"
)

func (n *node) string() string {
	return fmt.Sprintf("%s:%d (%s)", n.IP, n.ServicePort, n.Hostname)
}

func (n *node) save() {
	err := dbClient.Update(func(tx *bolt.Tx) error {
		jsonData, _ := json.Marshal(n)
		if err := tx.Bucket([]byte(nodeBucket)).Put([]byte(n.ID), jsonData); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatalln("Error updating node", err)
	}
}

func (n *node) delete() {
	err := dbClient.Update(func(tx *bolt.Tx) error {
		if err := tx.Bucket([]byte(nodeBucket)).Delete([]byte(n.ID)); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatalln("Error deleting node", err)
	}
}

func (n *node) getNodeClient() service.NodeClient {
	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", n.IP, n.ServicePort), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	// Creates a new CustomerClient
	n.ServiceClient = service.NewNodeClient(conn)
	return n.ServiceClient
}

func (n *node) execute(command string) map[string]*service.ExecutionResponse {
	response := make(map[string]*service.ExecutionResponse, 0)

	nodeResponse, err := n.getNodeClient().Execute(context.Background(), &service.ExecutionRequest{Command: command})
	if err != nil {
		log.Fatal(err)
	} else {
		response[n.Nickname] = nodeResponse
		n.updateLastSeen()
	}

	return response
}

func (n *node) details() *service.NodeDetails {
	response, err := n.getNodeClient().Details(context.Background(), &service.Empty{})
	if err != nil {
		log.Fatal(err)
	}

	return response
}

func (n *node) ping() map[string]bool {
	response := make(map[string]bool, 0)

	_, err := n.getNodeClient().Ping(context.Background(), &service.Empty{})
	if err != nil {
		response[n.Nickname] = false
	} else {
		response[n.Nickname] = true
		n.updateLastSeen()
	}

	return response
}

func (n *node) updateLastSeen() {
	n.LastSeen = time.Now()
	n.save()
}

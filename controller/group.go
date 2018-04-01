package controller

import (
	"fmt"
	"log"
	"strings"

	"github.com/bah2830/cluster/service"
	"github.com/boltdb/bolt"
)

func (g *group) delete(dbClient *bolt.DB) {
	err := dbClient.Update(func(tx *bolt.Tx) error {
		if err := tx.Bucket([]byte(groupBucket)).Delete([]byte(g.ID)); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatalln("Error deleting group", err)
	}
}

func (g *group) execute(command string) map[string]*service.ExecutionResponse {
	response := make(map[string]*service.ExecutionResponse, 0)

	responseCh := make(chan map[string]*service.ExecutionResponse, 0)

	executions := 0
	for _, n := range g.Nodes {
		executions++
		go func(n *node, responseCh chan map[string]*service.ExecutionResponse) {
			nodeResponse := n.execute(command)
			responseCh <- nodeResponse
		}(n, responseCh)
	}

	responded := 0
	for {
		select {
		case nodeResponse := <-responseCh:
			responded++
			for nodeName, r := range nodeResponse {
				response[nodeName] = r
			}

			if responded == executions {
				return response
			}
		}

	}
}

func (g *group) getNode(identifier string) (*node, error) {
	for _, node := range g.Nodes {
		idParts := strings.Split(node.ID, "-")
		shortID := idParts[len(idParts)-1]

		if node.ID == identifier || shortID == identifier {
			return node, nil
		}

		if node.Nickname == identifier {
			return node, nil
		}

		if node.Hostname == identifier {
			return node, nil
		}
	}

	return nil, fmt.Errorf("No node with identifier %s found", identifier)
}

package node

import (
	"github.com/bah2830/cluster/mdns"
	"github.com/davecgh/go-spew/spew"
)

type node struct {
	controller string
}

func Start(port string) {
	go startControllerDiscovery()
	startServer(port)
}

func startControllerDiscovery() {
	entries := mdns.StartDiscovery()
	for entry := range entries {
		spew.Dump(entry)
	}
}

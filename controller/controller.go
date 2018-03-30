package controller

import (
	"strconv"

	"github.com/bah2830/cluster/mdns"
)

type controller struct {
	nodes string
}

func Start(port string) {
	portInt, _ := strconv.Atoi(port)
	mdns := mdns.StartServer(portInt)
	defer mdns.Shutdown()

	startServer(port)
}

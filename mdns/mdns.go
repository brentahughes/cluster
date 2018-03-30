package mdns

import (
	"os"

	"github.com/hashicorp/mdns"
)

const serviceName = "cluster.controller"

func StartServer(port int) *mdns.Server {
	host, _ := os.Hostname()
	service, _ := mdns.NewMDNSService(host, serviceName, "", "", port, nil, []string{})
	server, _ := mdns.NewServer(&mdns.Config{Zone: service})
	return server
}

func StartDiscovery() chan *mdns.ServiceEntry {
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	mdns.Lookup(serviceName, entriesCh)
	return entriesCh
}

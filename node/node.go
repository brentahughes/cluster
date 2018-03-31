package node

import (
	"log"
	"net"
	"os"
	"strings"

	"github.com/bah2830/cluster/service"
	"github.com/hashicorp/mdns"
	"github.com/spf13/viper"
)

type Node struct{}

type controller struct {
	ip       string
	port     int
	hostname string
}

func (n *Node) Start() {
	go startMDNS()
	n.startAPIServer()
}

func (n *Node) GetNodeDetails() *service.NodeDetails {
	host, _ := os.Hostname()

	nodeDetails := &service.NodeDetails{
		Hostname:    host,
		ServicePort: int32(viper.GetInt("rpc.port")),
		Version:     viper.GetString("version"),
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			panic(err)
		}

		network := &service.NodeNetworkDetails{Interface: i.Name}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip.IsLoopback() == false && strings.Contains(ip.String(), ".") {
				network.IpAddress = ip.String()
			}

			if i.HardwareAddr.String() != "" {
				network.MacAddress = i.HardwareAddr.String()
			}
		}

		if network.IpAddress != "" && network.MacAddress != "" {
			nodeDetails.Networks = append(nodeDetails.Networks, network)
		}
	}

	return nodeDetails
}

func startMDNS() *mdns.Server {
	log.Printf("Starting mdns service advertising on %s", viper.GetString("mdns.service"))
	host, _ := os.Hostname()

	mdnsService, _ := mdns.NewMDNSService(
		host,
		viper.GetString("mdns.service"),
		"",
		"",
		viper.GetInt("rpc.port"),
		nil,
		[]string{"cluster_management"},
	)

	server, _ := mdns.NewServer(&mdns.Config{Zone: mdnsService})
	return server
}

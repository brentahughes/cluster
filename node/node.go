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

type Node struct {
	Version string
}

func (n *Node) Start() {
	server := n.startMDNS()
	defer server.Shutdown()
	n.startAPIServer()
}

func (n *Node) GetNodeDetails() *service.NodeDetails {
	host, _ := os.Hostname()

	nodeDetails := &service.NodeDetails{
		Hostname:    host,
		ServicePort: int32(viper.GetInt("rpc.port")),
		Version:     n.Version,
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

func (n *Node) startMDNS() *mdns.Server {
	details := n.GetNodeDetails()
	ip := details.Networks[0].IpAddress
	addr := net.ParseIP(ip)

	log.Printf("Starting mdns service advertising on %s (%s)", ip, viper.GetString("mdns.service"))

	host, _ := os.Hostname()

	mdnsService, err := mdns.NewMDNSService(
		host,
		viper.GetString("mdns.service"),
		"",
		"",
		viper.GetInt("rpc.port"),
		[]net.IP{addr},
		[]string{"cluster_management"},
	)
	if err != nil {
		log.Fatal(err)
	}

	server, _ := mdns.NewServer(&mdns.Config{Zone: mdnsService})
	return server
}

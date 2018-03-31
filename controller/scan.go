package controller

import (
	"log"
	"time"

	"github.com/hashicorp/mdns"
	"github.com/spf13/viper"
)

func (c *Controller) FindNodes() {
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	defer close(entriesCh)

	mdns.Lookup(viper.GetString("mdns.service"), entriesCh)

	duration := time.After(viper.GetDuration("scan.duration"))
	nodes := make([]*node, 0)

LOOPBREAK:
	for {
		select {
		case entry := <-entriesCh:
			node := &node{
				hostname:    entry.Host,
				ip:          entry.AddrV4.String(),
				servicePort: entry.Port,
				lastSeen:    time.Now(),
			}

			log.Println("Found node " + node.string())
			nodes = append(nodes, node)
		case <-duration:
			break LOOPBREAK
		}
	}

	log.Printf("Found %d nodes", len(nodes))
	c.nodes = nodes
}

package controller

import (
	"log"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/hashicorp/mdns"
	"github.com/spf13/viper"
)

func (c *Controller) FindNodes() {
	log.Printf("Searcing for cluster nodes on mdns service %s", viper.GetString("mdns.service"))
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	defer close(entriesCh)

	mdns.Lookup(viper.GetString("mdns.service"), entriesCh)

	duration := time.After(viper.GetDuration("scan.duration"))
	nodes := make([]*node, 0)

LOOPBREAK:
	for {
		select {
		case entry := <-entriesCh:
			currentTime := time.Now()

			hostname := strings.TrimSuffix(entry.Host, ".")

			node := &node{
				ID:          uuid.New().String(),
				Nickname:    hostname,
				Hostname:    hostname,
				IP:          entry.AddrV4.String(),
				ServicePort: entry.Port,
				FirstSeen:   currentTime,
				LastSeen:    currentTime,
			}

			log.Println("Found node " + node.string())

			if found, err := c.nodeExists(node); found {
				log.Printf("Duplicate node %s found: %s", node.string(), err.Error())
			} else {
				nodes = append(nodes, node)
				c.nodeChanges = true
			}
		case <-duration:
			break LOOPBREAK
		}
	}

	log.Printf("Found %d new nodes", len(nodes))
	c.nodes = append(c.nodes, nodes...)
}

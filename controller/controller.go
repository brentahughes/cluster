package controller

import (
	"fmt"
	"time"

	"github.com/bah2830/cluster/service"
)

type Controller struct {
	nodes []*node
}

type node struct {
	hostname      string
	ip            string
	servicePort   int
	serviceClient *service.NodeClient
	details       *service.NodeDetails
	lastSeen      time.Time
}

func (n *node) string() string {
	return fmt.Sprintf("%s:%d (%s)", n.ip, n.servicePort, n.hostname)
}

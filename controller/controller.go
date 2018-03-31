package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"
	"text/tabwriter"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/bah2830/cluster/service"
	"github.com/boltdb/bolt"
	"github.com/google/uuid"
	"github.com/hashicorp/mdns"
	"github.com/spf13/viper"
)

const (
	nodeBucket  = "nodes"
	groupBucket = "groups"
)

type Controller struct {
	nodes       []*node
	nodeChanges bool
	dbClient    *bolt.DB
}

type node struct {
	ID            string
	Nickname      string
	Hostname      string
	IP            string
	ServicePort   int
	ServiceClient service.NodeClient
	LastSeen      time.Time
	FirstSeen     time.Time
}

func GetController() *Controller {
	user, _ := user.Current()

	// Init database directory if it doesn't already exist
	if _, err := os.Stat(user.HomeDir + "/.cluster"); os.IsNotExist(err) {
		log.Println("Initializing ~/.cluster directory")
		os.Mkdir(user.HomeDir+"/.cluster", 0700)
	}

	c := &Controller{}
	c.getDBConnection()
	c.loadNodes()
	c.loadGroups()

	return c
}

func (c *Controller) CleanExit() {
	c.saveNodes()
	c.saveGroups()

	if c.dbClient != nil {
		c.dbClient.Close()
	}
}

func (c *Controller) getDBConnection() *bolt.DB {
	if c.dbClient != nil {
		return c.dbClient
	}

	db, err := bolt.Open(viper.GetString("cluster.db"), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		for _, bucket := range []string{nodeBucket, groupBucket} {
			_, err := tx.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				return fmt.Errorf("Error setting up bucket: %s. %s", bucket, err)
			}
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	c.dbClient = db

	return db

}

func (c *Controller) loadNodes() {
	nodes := make([]*node, 0)
	c.dbClient.View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket([]byte(nodeBucket)).Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			node := &node{}
			json.Unmarshal(v, node)
			nodes = append(nodes, node)
		}

		return nil
	})

	c.nodes = nodes
}

func (c *Controller) saveNodes() {
	if c.nodeChanges == false {
		return
	}

	c.dbClient.Update(func(tx *bolt.Tx) error {
		for _, node := range c.nodes {
			nodeCopy := *node
			nodeCopy.ServiceClient = nil
			jsonData, _ := json.Marshal(nodeCopy)
			if err := tx.Bucket([]byte(nodeBucket)).Put([]byte(node.ID), jsonData); err != nil {
				return err
			}
		}

		return nil
	})
}

func (c *Controller) loadGroups() {

}

func (c *Controller) saveGroups() {

}

func (c *Controller) ListNodes() {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 4, '\t', 0)

	fmt.Fprintf(w, "Name\tHostname\tIP:Port\tLast Seen\tID\n")
	for _, node := range c.nodes {
		idParts := strings.Split(node.ID, "-")
		shortID := idParts[len(idParts)-1]

		fmt.Fprintf(
			w,
			"%s\t%s\t%s:%d\t%s\t%s\n",
			node.Nickname,
			node.Hostname,
			node.IP,
			node.ServicePort,
			node.LastSeen.Format("2006-01-02 15:04:05"),
			shortID,
		)
	}

	w.Flush()
}

func (c *Controller) nodeExists(n *node) (bool, error) {
	for _, node := range c.nodes {
		if node.IP == n.IP {
			return true, errors.New("Duplicate ip")
		}

		if node.Hostname == n.Hostname {
			return true, errors.New("Duplicate hostname")
		}

		if node.Nickname == n.Nickname {
			return true, errors.New("Duplicate name")
		}
	}

	return false, nil
}

func (c *Controller) SetNodeNickName(identifier, name string) {
	node, err := c.getNode(identifier)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := c.getNode(name); err == nil {
		log.Fatalf("Node with identifier %s already exists", name)
	}

	node.Nickname = name
	node.save(c.dbClient)
}

func (c *Controller) DeleteNode(identifier string) {
	n, err := c.getNode(identifier)
	if err != nil {
		log.Fatal(err)
	}
	n.delete(c.dbClient)
}

func (c *Controller) getNode(identifier string) (*node, error) {
	for _, node := range c.nodes {
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

func (c *Controller) Execute(identifier, command string) {
	n, err := c.getNode(identifier)
	if err != nil {
		log.Fatal(err)
	}

	response := n.execute(command)
	fmt.Printf(response.StdOut)
}

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

func (c *Controller) NodeDetails(identifier string) {
	node, err := c.getNode(identifier)
	if err != nil {
		log.Fatal(err)
	}

	details := node.details()
	jsonString, _ := yaml.Marshal(details)
	fmt.Println(string(jsonString))
}

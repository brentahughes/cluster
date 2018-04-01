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
	nodes        []*node
	nodeChanges  bool
	groups       []*group
	groupChanges bool
	dbClient     *bolt.DB
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

type group struct {
	ID         string
	Nickname   string
	Nodes      []*node
	NodeNames  []string
	CreateDate time.Time
}

type executer interface {
	execute(string) map[string]*service.ExecutionResponse
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
	groups := make([]*group, 0)
	c.dbClient.View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket([]byte(groupBucket)).Cursor()

		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			group := &group{}
			json.Unmarshal(v, group)

			for _, nodeName := range group.NodeNames {
				node, err := c.getNode(nodeName)
				if err != nil {
					log.Println("WARNING: Unable to load all nodes for group", err)
				} else {
					group.Nodes = append(group.Nodes, node)
				}
			}

			groups = append(groups, group)
		}

		return nil
	})

	c.groups = groups
}

func (c *Controller) saveGroups() {
	if c.groupChanges == false {
		return
	}

	c.dbClient.Update(func(tx *bolt.Tx) error {
		for _, group := range c.groups {
			groupCopy := *group
			groupCopy.Nodes = make([]*node, 0)
			jsonData, _ := json.Marshal(groupCopy)
			if err := tx.Bucket([]byte(groupBucket)).Put([]byte(group.ID), jsonData); err != nil {
				return err
			}
		}

		return nil
	})
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

func (c *Controller) ListGroups() {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 4, '\t', 0)

	fmt.Fprintf(w, "Name\tNodes\tCreated\tID\n")
	for _, group := range c.groups {
		idParts := strings.Split(group.ID, "-")
		shortID := idParts[len(idParts)-1]

		fmt.Fprintf(
			w,
			"%s\t%d\t%s\t%s\n",
			group.Nickname,
			len(group.Nodes),
			group.CreateDate.Format("2006-01-02 15:04:05"),
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

func (c *Controller) groupExists(g *group) (bool, error) {
	for _, group := range c.groups {
		if group.Nickname == g.Nickname {
			return true, errors.New("Duplicate group name")
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

	if _, err := c.getGroup(name); err == nil {
		log.Fatalf("Group with identifier %s already exists", name)
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

func (c *Controller) getGroup(identifier string) (*group, error) {
	for _, group := range c.groups {
		idParts := strings.Split(group.ID, "-")
		shortID := idParts[len(idParts)-1]

		if group.ID == identifier || shortID == identifier {
			return group, nil
		}

		if group.Nickname == identifier {
			return group, nil
		}
	}

	return nil, fmt.Errorf("No group with identifier %s found", identifier)
}

func (c *Controller) Execute(identifier, command string) {
	var n executer
	var err error

	n, err = c.getNode(identifier)
	if err != nil {
		// Check if it is the name of a group
		n, err = c.getGroup(identifier)
		if err != nil {
			log.Fatal("No group or node found with the identifier", identifier)
		}
	}

	responses := n.execute(command)

	for nodeName, response := range responses {
		fmt.Println("Response from", nodeName)
		fmt.Println("====================================================================")
		fmt.Println(response.StdOut)
		fmt.Println(response.StdErr)
	}
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

func (c *Controller) CreateGroup(name string, nodes []string) {
	group := &group{
		ID:         uuid.New().String(),
		Nickname:   name,
		NodeNames:  nodes,
		CreateDate: time.Now(),
	}

	if exists, err := c.groupExists(group); exists {
		log.Fatal(err)
	}

	if exists, err := c.nodeExists(&node{Nickname: name}); exists {
		log.Fatal(err)
	}

	for _, nodeName := range nodes {
		if _, err := c.getNode(nodeName); err != nil {
			log.Fatal(err)
		}
	}
	c.groups = append(c.groups, group)
	c.groupChanges = true
}

func (c *Controller) DeleteGroup(identifier string) {
	g, err := c.getGroup(identifier)
	if err != nil {
		log.Fatal(err)
	}

	g.delete(c.dbClient)
}

func (c *Controller) AddNodesToGroup(identifer string, nodes []string) {
	group, err := c.getGroup(identifer)
	if err != nil {
		log.Fatal(err)
	}

	for _, nodeName := range nodes {
		node, _ := group.getNode(nodeName)
		if node != nil {
			log.Fatalf("Node %s is already a part of this group", nodeName)
		}

		node, err = c.getNode(nodeName)
		if err != nil {
			log.Fatal(err)
		}

		group.NodeNames = append(group.NodeNames, nodeName)
		group.Nodes = append(group.Nodes, node)
		c.groupChanges = true
	}
}

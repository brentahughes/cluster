# Cluster

Controller node cluster management using gRPC and mDNS service discovery. Each node will watch for mDNS entries on the network. When the controller entry is found it will send a checkin message with it's details to the controllers gRPC service. From then on the controller and communicate with all nodes that have checked in.

## Usage
```bash
Cluster managament application

Usage:
  cluster [command]

Available Commands:
  execute     execute a command on a node
  group       Commands for managing groups withing the cluster
  help        Help about any command
  node        Commands for managing nodes withing the cluster
  ping        Send a ping to each node

Flags:
  -c, --config string   Path to cluster config.db file (default "/Users/username/.cluster/config.db")
  -h, --help            help for cluster
      --version         version for cluster

Use "cluster [command] --help" for more information about a command.
```

### Node
#### Deploy Node
```bash
./cluster node deploy
```

#### Find Nodes
```bash
./cluster node scan
```

#### Rename Node
```bash
./cluster node name node1 new_node1
```

#### List Nodes
```bash
./cluster node list

Name		Hostname	    IP:Port			        Last Seen		        ID
node1		bdd4e94392d2	192.168.1.135:10000	    2018-04-01 13:42:14	    f4f18505b85b
node2		000107e91376	192.168.1.137:10000	    2018-04-01 20:51:14	    3aeac64e3231
node3		8c1a750161ac	192.168.1.150:10000	    2018-04-01 14:04:53	    7a971aec1431
node4		raspberrypi	192.168.1.128:10000	    2018-04-01 13:58:17	    8c8f0945ce28
```

#### Node Online Check
```bash
./cluster ping node1

NODE		ONLINE
node1		true
```

### Group

#### Create Group
```bash
./cluster group create group1 node1 node2 node3 node4
```

#### Create Group
```bash
./cluster group create group1 node1 node2 node3 node4
```

#### List Groups
```bash
./cluster group list

Name		Nodes		Created			        ID
group		4		    2018-04-01 20:27:17	    7bb3f4a7ea8e
```

#### Group Details
```bash
./cluster group details group1

id: 47282399-3344-4815-916d-7bb3f4a7ea8e
nickname: group1
nodes:
- id: 1312e7bf-b4ff-4ad3-8e70-f4f18505b85b
  nickname: node1
  hostname: bdd4e94392d2
  ip: 192.168.1.135
  serviceport: 10000
  serviceclient: null
  lastseen: 2018-04-01T21:15:29.178787-05:00
  firstseen: 2018-04-01T13:42:14.754074-05:00
- id: 2f13d3cc-12ad-468b-8e97-3aeac64e3231
  nickname: node2
  hostname: 000107e91376
  ip: 192.168.1.137
  serviceport: 10000
  serviceclient: null
  lastseen: 2018-04-01T20:51:14.640109-05:00
  firstseen: 2018-04-01T14:04:03.788414-05:00
- id: 3d1bebe8-34e3-416e-9c0d-7a971aec1431
  nickname: node3
  hostname: 8c1a750161ac
  ip: 192.168.1.150
  serviceport: 10000
  serviceclient: null
  lastseen: 2018-04-01T14:04:53.525019-05:00
  firstseen: 2018-04-01T14:04:53.525019-05:00
- id: b36fa5c7-91b1-4802-b508-8c8f0945ce28
  nickname: node4
  hostname: raspberrypi
  ip: 192.168.1.128
  serviceport: 10000
  serviceclient: null
  lastseen: 2018-04-01T13:58:17.346852-05:00
  firstseen: 2018-04-01T13:58:17.346852-05:00
nodenames:
- node1
- node2
- node3
- node4
createdate: 2018-04-01T20:27:17.824959-05:00
```

#### Execute Command On Group
```bash
./cluster execute group1 -- lsb_release -da

Response from node2
====================================================================
Distributor ID:	Raspbian
Description:	Raspbian GNU/Linux 9.1 (stretch)
Release:	9.1
Codename:	stretch


Response from node1
====================================================================
Distributor ID:	Raspbian
Description:	Raspbian GNU/Linux 9.1 (stretch)
Release:	9.1
Codename:	stretch


Response from node3
====================================================================
Distributor ID:	Raspbian
Description:	Raspbian GNU/Linux 9.1 (stretch)
Release:	9.1
Codename:	stretch


Response from node4
====================================================================
Distributor ID:	Raspbian
Description:	Raspbian GNU/Linux 9.1 (stretch)
Release:	9.1
Codename:	stretch
```
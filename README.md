# Cluster

Controller node cluster management using gRPC and mDNS service discovery. Each node will watch for mDNS entries on the network. When the controller entry is found it will send a checkin message with it's details to the controllers gRPC service. From then on the controller and communicate with all nodes that have checked in.
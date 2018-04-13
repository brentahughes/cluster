package main

//go:generate protoc -I service/ service/service.proto --go_out=plugins=grpc:service

import "github.com/bah2830/cluster/cmd"

const VERSION = "0.2.0"

func main() {
	cmd.Execute(VERSION)
}

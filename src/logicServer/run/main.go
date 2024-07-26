package main

import (
	"myGameDemo/logicServer"
)

func main() {
	server := logicServer.LogicServer{
		Addr: ":8080",
	}
	server.Run()
}

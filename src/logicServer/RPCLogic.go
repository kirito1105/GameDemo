package logicServer

import (
	"google.golang.org/grpc"
	"log"
	"myGameDemo/myRPC"
	"sync"
)

type RPCLogic struct{}

var logicRPC *RPCLogic
var once sync.Once

func GetLogicRPC() *RPCLogic {
	once.Do(func() {
		logicRPC = &RPCLogic{}
	})
	return logicRPC
}

// 获取连接rc的client
var rcClient myRPC.RCenterRPCClient
var once1 sync.Once

func GetRcClient() myRPC.RCenterRPCClient {
	once1.Do(func() {
		conn, err := grpc.Dial("localhost:25565", grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}
		rcClient = myRPC.NewRCenterRPCClient(conn)
	})
	return rcClient
}

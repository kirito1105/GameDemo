package rcenterServer

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"myGameDemo/myRPC"
	"net"
	"sync"
)

type RCServer struct {
	*myRPC.UnimplementedRCenterRPCServer
}

func (p *RCServer) CreateRoom(ctx context.Context, gameFindInfo *myRPC.GameRoomFindInfo) (*myRPC.RoomInfo, error) {

	room, err := CreateRoom(gameFindInfo)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (p *RCServer) RoomServerHeart(ctx context.Context, info *myRPC.RoomServerInfo) (*myRPC.Res, error) {
	//fmt.Println(info)
	RoomServerHeart(info)
	//fmt.Println(GetRoomServerRegisterCenter().roomServerList.Front())
	return &myRPC.Res{Code: myRPC.Code_SUCCESS}, nil
}

type RPCRoomCenter struct{}

var roomCenterRPC *RPCRoomCenter
var once sync.Once

func GetLogicRPC() *RPCRoomCenter {
	once.Do(func() {
		roomCenterRPC = &RPCRoomCenter{}
	})
	return roomCenterRPC
}

func (p *RPCRoomCenter) run() {
	grpcServer := grpc.NewServer()
	myRPC.RegisterRCenterRPCServer(grpcServer, new(RCServer))

	lis, err := net.Listen("tcp", ":25565")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer.Serve(lis)
}
func (p *RPCRoomCenter) server() {
	p.run()
}

//调用gpc接口

func CreateRoomClient(addr string) myRPC.RoomRPCClient {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	return myRPC.NewRoomRPCClient(conn)
}

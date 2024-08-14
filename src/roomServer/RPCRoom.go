package roomServer

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log"
	"myGameDemo/myRPC"
	"net"
	"strconv"
	"sync"
)

type RServer struct {
	*myRPC.UnimplementedRoomRPCServer
}

func (R *RServer) CreateRoom(ctx context.Context, info *myRPC.GameRoomFindInfo) (*myRPC.RoomInfo, error) {
	roominfo := GetRoomController().RoomCreate()
	return roominfo, nil
}
func (R *RServer) FindARoom(ctx context.Context, info *myRPC.GameRoomFindInfo) (*myRPC.RoomInfo, error) {
	list := GetRoomController().FindRooms(info, 1)
	if len(list) > 0 {
		return list[0], nil
	}
	return GetRoomController().RoomCreate(), nil
}

type RPCRoom struct {
	ip   string
	port int
}

var roomRPC *RPCRoom
var once sync.Once

var (
	ip   string
	port int
)

func SetAddr(ip1 string, port1 int) {
	ip = ip1
	port = port1
}

func GetRoomRPC() *RPCRoom {
	once.Do(func() {
		roomRPC = &RPCRoom{ip, port}
	})
	return roomRPC
}

func (p *RPCRoom) run() {
	grpcServer := grpc.NewServer()
	myRPC.RegisterRoomRPCServer(grpcServer, new(RServer))

	lis, err := net.Listen("tcp", p.ip+":"+strconv.Itoa(p.port))
	if err != nil {
		logrus.Fatal("[RPC]输出化失败", err)
	}
	logrus.Info("[RPC]初始化成功")
	err = grpcServer.Serve(lis)
	if err != nil {
		logrus.Fatal("[RPC]调用出错", err)
		return
	}
}

func (p *RPCRoom) server() {
	p.run()

}

var myClient myRPC.RCenterRPCClient
var once1 sync.Once

func GetMyClient() myRPC.RCenterRPCClient {
	once1.Do(func() {
		conn, err := grpc.Dial("localhost:25565", grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}
		myClient = myRPC.NewRCenterRPCClient(conn)
	})
	return myClient
}

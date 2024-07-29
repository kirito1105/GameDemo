package roomServer

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"myGameDemo/myRPC"
	"net"
	"strconv"
	"sync"
	"time"
)

type RServer struct {
	*myRPC.UnimplementedRoomRPCServer
}

func RoomCreate() *myRPC.RoomInfo {
	theRoom := NewRoom()
	theRoom.Start()
	roomId := strconv.Itoa(int(time.Now().Unix()))
	theRoom.RoomID = roomId

	GetRoomController().AddRoom(theRoom, roomId)
	roominfo := &myRPC.RoomInfo{
		IsFind:   true,
		RoomId:   roomId,
		RoomAddr: theRoom.GetTCPAddr().String(),
	}

	fmt.Println(GetRoomController().Summary())
	return roominfo
}

func (R RServer) CreateRoom(ctx context.Context, info *myRPC.GameRoomFindInfo) (*myRPC.RoomInfo, error) {
	roominfo := RoomCreate()
	return roominfo, nil
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
		log.Fatal(err)
	}

	grpcServer.Serve(lis)
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

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

func RoomStart() *myRPC.RoomInfo {
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	cli, _ := net.ListenTCP("tcp", addr)
	//TODO room信息处理
	roomId := strconv.Itoa(int(time.Now().Unix()))
	GetRoomController().AddRoom(Room{
		Addr:          cli.Addr().String(),
		OnlinePlayers: make([]player, 0),
	}, roomId)

	roominfo := &myRPC.RoomInfo{
		IsFind:   true,
		RoomId:   roomId,
		RoomAddr: cli.Addr().String(),
	}
	//TODO 抽象房间为类
	go func() {
		for {
			//TODO 房间运行
			id := roomId
			conn, err := cli.Accept()
			if err != nil {
				//TODO TCP出错
				return
			}
			//TODO 处理玩家信息
			//Todo 获取玩家ID
			username := "123"
			GetRoomController().PlayerOnline(player{username: username}, id)
			//TODO 将与玩家的连接加入连接总表
			fmt.Println(conn)
		}

	}()
	fmt.Println(GetRoomController().Summary())
	return roominfo
}

func (R RServer) CreateRoom(ctx context.Context, info *myRPC.GameRoomFindInfo) (*myRPC.RoomInfo, error) {
	roominfo := RoomStart()
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

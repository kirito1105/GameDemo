package rcenterServer

import (
	"context"
	"fmt"
	list "github.com/liyue201/gostl/ds/list/bidlist"
	"myGameDemo/myRPC"
	"sync"
)

type RoomServerNode struct {
	Addr      string              //服务器grpc地址
	PlayerNum int                 //当前服务器在线玩家
	RoomNum   int                 //当前服务器正在进行的对局
	Client    myRPC.RoomRPCClient //服务器的rpc Client
}

type RoomServerRegisterCenter struct {
	roomServerList list.List[RoomServerNode]
	mutex          sync.Mutex
}

func (rc *RoomServerRegisterCenter) RegNewServer(info *myRPC.RoomServerInfo) {

	rc.mutex.Lock()
	defer rc.mutex.Unlock()
	fmt.Println(info)
	for n := rc.roomServerList.FrontNode(); n != nil; n = n.Next() {
		if n.Value.Addr == info.Addr {
			n.Value.RoomNum = int(info.RoomNum)
			n.Value.PlayerNum = int(info.PlayerNum)
			return
		}
	}
	tmp := RoomServerNode{
		Addr:      info.Addr,
		PlayerNum: int(info.PlayerNum),
		RoomNum:   int(info.RoomNum),
		Client:    CreateRoomClient(info.Addr),
	}
	rc.roomServerList.PushBack(tmp)
}

func (rc *RoomServerRegisterCenter) minRoomServe() myRPC.RoomRPCClient {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()
	if rc.roomServerList.Len() == 0 {
		return nil
	}
	minNum := rc.roomServerList.Front().RoomNum
	cl := rc.roomServerList.Front().Client

	for n := rc.roomServerList.FrontNode(); n != nil; n = n.Next() {
		if minNum > n.Value.RoomNum {
			minNum = n.Value.RoomNum
			cl = n.Value.Client
		}
	}
	return cl
}

func (rc *RoomServerRegisterCenter) GetRoomList() (*myRPC.RoomInfoArray, error) {
	arr := &myRPC.RoomInfoArray{
		Rooms: make([]*myRPC.RoomInfoNode, 0),
	}
	for n := rc.roomServerList.FrontNode(); n != nil; n = n.Next() {
		list, err := n.Value.Client.GetRoomList(context.Background(), &myRPC.Empty{})
		if err != nil {
			continue
		}
		for _, i := range list.Rooms {
			i.Addr = n.Value.Addr
		}
		arr.Rooms = append(arr.Rooms, list.Rooms...)
		if len(arr.Rooms) > 20 {
			return arr, nil
		}
	}
	return arr, nil
}
func (rc *RoomServerRegisterCenter) FindRoomWithAddr(addr string) myRPC.RoomRPCClient {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()
	for i := rc.roomServerList.FrontNode(); i != nil; i = i.Next() {
		if i.Value.Addr == addr {
			return i.Value.Client
		}
	}
	return nil
}

var roomServerRegisterCenter *RoomServerRegisterCenter
var once1 sync.Once

func GetRoomServerRegisterCenter() *RoomServerRegisterCenter {
	once1.Do(func() {
		roomServerRegisterCenter = &RoomServerRegisterCenter{}
	})
	return roomServerRegisterCenter
}

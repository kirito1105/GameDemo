package rcenterServer

import (
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

func (rc *RoomServerRegisterCenter) minPlayerServe() myRPC.RoomRPCClient {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()
	if rc.roomServerList.Len() == 0 {
		return nil
	}
	minNum := rc.roomServerList.Front().PlayerNum
	cl := rc.roomServerList.Front().Client

	for n := rc.roomServerList.FrontNode(); n != nil; n = n.Next() {
		if minNum > n.Value.PlayerNum {
			minNum = n.Value.PlayerNum
			cl = n.Value.Client
		}
	}

	return cl
}

var roomServerRegisterCenter *RoomServerRegisterCenter
var once1 sync.Once

func GetRoomServerRegisterCenter() *RoomServerRegisterCenter {
	once1.Do(func() {
		roomServerRegisterCenter = &RoomServerRegisterCenter{}
	})
	return roomServerRegisterCenter
}

package roomServer

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"myGameDemo/myRPC"
	"strconv"
	"sync"
	"time"
)

const MAX_NUM = 4

type playerInfoSum struct {
	username string
}

func (p *playerInfoSum) String() string {
	return p.username
}

type RoomInfoSum struct {
	Addr          string
	OnlinePlayers []playerInfoSum
}

type RoomController struct {
	RoomwithId map[string]*Room
	players    map[string]*PlayerTask
	mutex      sync.RWMutex
}

var roomController *RoomController
var once2 sync.Once

func GetRoomController() *RoomController {
	once2.Do(func() {
		roomController = &RoomController{
			RoomwithId: make(map[string]*Room),
			players:    make(map[string]*PlayerTask),
		}
	})
	return roomController
}

func (rc *RoomController) AddRoom(room *Room, id string) {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()
	rc.RoomwithId[id] = room
}

// DeleteSlice 删除指定元素。
func DeleteSliceR(s []*Room, roomid string) []*Room {
	j := 0
	for _, v := range s {
		if v.RoomID != roomid {
			s[j] = v
			j++
		}
	}
	return s[:j]
}

func (rc *RoomController) RemoveRoom(id string) {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()
	delete(rc.RoomwithId, id)
}

func (rc *RoomController) GetRoom(id string) *Room {
	rc.mutex.RLock()
	defer rc.mutex.RUnlock()
	if val, ok := rc.RoomwithId[id]; ok {
		return val
	} else {
		return nil
	}
}

func (rc *RoomController) PlayerOffline(username string) {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()
	delete(rc.players, username)
	//todo
}

func (rc *RoomController) PlayerOnline(player *PlayerTask) {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()
	rc.players[player.username] = player
	//todo
}

func (rc *RoomController) Summary() string {
	rc.mutex.RLock()
	defer rc.mutex.RUnlock()
	var sum string
	for key, v := range rc.RoomwithId {
		sum += "Roomid:" + key + " " + "Addr:" + v.GetTCPAddr().String()

		sum += "\n"
	}
	return sum
}

func (rc *RoomController) RoomCreate() *myRPC.RoomInfo {
	theRoom := NewRoom()
	theRoom.Start()

	roomId := strconv.Itoa(int(time.Now().UnixNano()))
	theRoom.RoomID = roomId
	logrus.Debug("[房间]房间创建成功:", theRoom.RoomID)
	rc.AddRoom(theRoom, roomId)
	roominfo := &myRPC.RoomInfo{
		IsFind:   true,
		RoomId:   roomId,
		RoomAddr: theRoom.GetTCPAddr().String(),
	}

	fmt.Println(rc.Summary())
	return roominfo
}

func (this *RoomController) FindRooms(info *myRPC.GameRoomFindInfo, num int) []*myRPC.RoomInfo {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	n := 0
	list := []*myRPC.RoomInfo{}
	for _, v := range this.RoomwithId {
		//todo
		in := &myRPC.RoomInfo{
			IsFind:   true,
			RoomId:   v.RoomID,
			RoomAddr: v.GetTCPAddr().String(),
		}
		list = append(list, in)
		n++
		if n == num {
			return list
		}
	}
	return list

}

func (rc *RoomController) GetRoomList() *myRPC.RoomInfoArray {
	arr := &myRPC.RoomInfoArray{
		Rooms: make([]*myRPC.RoomInfoNode, 0),
	}
	n := 0
	for id, room := range rc.RoomwithId {
		info := &myRPC.RoomInfoNode{
			RoomId:    id,
			PlayerNum: room.GetOnlineNum(),
		}
		if info.PlayerNum >= MAX_NUM {
			continue
		}
		arr.Rooms = append(arr.Rooms, info)
		n++
		if n > 20 {
			break
		}
	}
	return arr
}

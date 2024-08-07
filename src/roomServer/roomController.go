package roomServer

import (
	"sync"
)

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
	return rc.RoomwithId[id]
}

func (rc *RoomController) PlayerOffline(username string, id string) {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()
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

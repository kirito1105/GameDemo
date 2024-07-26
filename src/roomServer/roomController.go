package roomServer

import (
	"sync"
)

type player struct {
	username string
}

func (p *player) String() string {
	return p.username
}

type Room struct {
	Addr          string
	OnlinePlayers []player
}

type RoomController struct {
	Rooms map[string]Room
	mutex sync.RWMutex
}

var roomController *RoomController
var once2 sync.Once

func GetRoomController() *RoomController {
	once2.Do(func() {
		roomController = &RoomController{
			Rooms: make(map[string]Room),
		}
	})
	return roomController
}

func (rc *RoomController) AddRoom(room Room, id string) {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()
	rc.Rooms[id] = room
}

func (rc *RoomController) RemoveRoom(id string) {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()
	delete(rc.Rooms, id)
}

func (rc *RoomController) GetRoom(id string) Room {
	rc.mutex.RLock()
	defer rc.mutex.RUnlock()
	return rc.Rooms[id]
}

// DeleteSlice 删除指定元素。
func DeleteSlice(s []player, unsername string) []player {
	j := 0
	for _, v := range s {
		if v.username != unsername {
			s[j] = v
			j++
		}
	}
	return s[:j]
}

func (rc *RoomController) PlayerOffline(username string, id string) {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()
	tmp := rc.Rooms[id]
	DeleteSlice(tmp.OnlinePlayers, username)
	rc.Rooms[id] = tmp
}

func (rc *RoomController) PlayerOnline(user player, id string) {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()
	tmp := rc.Rooms[id]
	tmp.OnlinePlayers = append(tmp.OnlinePlayers, user)
	rc.Rooms[id] = tmp
}

func (rc *RoomController) Summary() string {
	rc.mutex.RLock()
	defer rc.mutex.RUnlock()
	var sum string
	for key, v := range rc.Rooms {
		sum += "Roomid:" + key + " " + "Addr:" + v.Addr
		for _, player := range v.OnlinePlayers {
			sum += " " + player.String()
		}
		sum += "\n"
	}
	return sum
}

package roomServer

import (
	"myGameDemo/mynet"
	world "myGameDemo/world"
	"net"
)

type Room struct {
	RoomID         string
	tcpListener    mynet.TCPListener
	players        []*PlayerTask
	playerWithName map[string]*PlayerTask
	world0         *world.World
}

func NewRoom() *Room {
	t := &Room{
		tcpListener:    mynet.TCPListener{},
		playerWithName: make(map[string]*PlayerTask),
		world0:         nil,
	}
	go func() {
		a := world.NewWorld()
		t.world0 = a
	}()
	return t
}

func (r *Room) GetTCPAddr() net.Addr {
	return r.tcpListener.Addr()
}

func (r *Room) GetWorld() *world.World {
	return r.world0
}

func (r *Room) Start() {
	err := r.tcpListener.ListenTCP()
	if err != nil {
		return
	}
	go func() {
		for {
			conn, err := r.tcpListener.AcceptTCP()
			if err != nil {
				return
			}
			p := NewPlayerTask(conn, r)
			r.players = append(r.players, p)
			p.Start()
		}
	}()
}

func (r *Room) GetMyWorld(x int, y int) *world.Block {
	//todo
	b := r.world0.GetBlock(x, y)
	return &b
}

func (r *Room) GetOnlineList() {
	//todo
	return
}

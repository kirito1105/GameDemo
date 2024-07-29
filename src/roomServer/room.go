package roomServer

import (
	"myGameDemo/mynet"
	"net"
)

type Room struct {
	RoomID      string
	tcpListener mynet.TCPListener
}

func NewRoom() *Room {
	return &Room{
		tcpListener: mynet.TCPListener{},
	}
}

func (r *Room) GetTCPAddr() net.Addr {
	return r.tcpListener.Addr()
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
			NewPlayerTask(conn, r).Start()
		}
	}()
}

func (r *Room) GetOnlineList() {
	return
}

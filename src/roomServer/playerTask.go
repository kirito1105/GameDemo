package roomServer

import (
	"fmt"
	"myGameDemo/mynet"
	"net"
)

type playerTask struct {
	tcpTask *mynet.TCPTask
	inRoom  *Room
}

func (p *playerTask) ParseMsg(data []byte) bool {
	fmt.Println("receive msg:", string(data))
	//TODO implement
	p.tcpTask.SendMsg([]byte("res TEST"))
	return true
}

func NewPlayerTask(conn *net.TCPConn, r *Room) *playerTask {
	t := &playerTask{
		tcpTask: mynet.NewTCPTask(conn),
		inRoom:  r,
	}
	t.tcpTask.Task = t
	return t
}

func (this *playerTask) Start() {
	this.tcpTask.Start()
}

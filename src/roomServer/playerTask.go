package roomServer

import (
	"fmt"
	"myGameDemo/mynet"
	"net"
)

type PlayerTask struct {
	tcpTask *mynet.TCPTask
	inRoom  *Room
}

func (p *PlayerTask) ParseMsg(data []byte) bool {
	fmt.Println("receive msg:", string(data))
	//TODO implement
	return true
}

func NewPlayerTask(conn *net.TCPConn, r *Room) *PlayerTask {
	t := &PlayerTask{
		tcpTask: mynet.NewTCPTask(conn),
		inRoom:  r,
	}
	t.tcpTask.Task = t
	return t
}

func (this *PlayerTask) Start() {
	this.tcpTask.Start()
}

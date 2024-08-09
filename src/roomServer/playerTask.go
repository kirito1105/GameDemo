package roomServer

import (
	"encoding/base64"
	"fmt"
	"google.golang.org/protobuf/proto"
	"myGameDemo/myMsg"
	"myGameDemo/mynet"
	"net"
)

type PlayerTask struct {
	username  string
	isChecked bool
	tcpTask   *mynet.TCPTask
	inRoom    *Room
}

func (p *PlayerTask) ParseMsg(data []byte) bool {
	//fmt.Println("receive msg:", data)
	//fmt.Println("receive msg:", string(data))
	var msg myMsg.MsgFromClient
	_ = proto.Unmarshal(data, &msg)
	//fmt.Println(msg.GetAuthentication().GetUsername())

	//验证身份
	if !p.isChecked && msg.GetCmd() == myMsg.Cmd_Authentication {
		decoded, _ := base64.StdEncoding.DecodeString(msg.GetAuthentication().GetToken())
		fmt.Println(msg)
		a := CheckToken(decoded, msg.GetAuthentication().GetUsername(), msg.GetAuthentication().GetAddr(), msg.GetAuthentication().GetRoomId())
		if a {
			p.username = msg.GetAuthentication().GetUsername()
			p.isChecked = true

			p.inRoom.PlayerIn(p)
			fmt.Println(msg.GetAuthentication().GetUsername(), "验证成功")

		} else {
			fmt.Println(1)
			return true
		}

		scene := p.inRoom.GetInitInfo(p.username)

		m := myMsg.MsgFromService{
			Scene: scene,
		}
		bytes, _ := proto.Marshal(&m)
		p.tcpTask.SendMsg(AddHeader(bytes))
	}

	//未验证的其他信息不做处理
	if !p.isChecked {
		return false
	}

	//移动
	if msg.GetCmd() == myMsg.Cmd_Move {
		move := &PlayerMove{
			username: p.username,
			velocity: &Vector2{
				x: msg.GetMove().GetX(),
				y: msg.GetMove().GetY(),
			},
		}
		fmt.Println(msg.Move)
		p.inRoom.chan_PlayerMove <- move
	}

	//停止移动
	if msg.GetCmd() == myMsg.Cmd_StopMove {
		move := &PlayerMove{
			username: p.username,
			velocity: nil,
		}
		p.inRoom.chan_PlayerMove <- move
	}

	//TODO implement
	return true
}

func NewPlayerTask(conn *net.TCPConn, r *Room) *PlayerTask {
	t := &PlayerTask{
		tcpTask:   mynet.NewTCPTask(conn),
		inRoom:    r,
		isChecked: false,
	}
	t.tcpTask.Task = t
	return t
}

func (this *PlayerTask) Start() {
	this.tcpTask.Start()
}

// 玩家局内信息
type Player struct {
	ObjBase
	Asta   myMsg.AnimatorStatus
	Online bool
}

func (this *Player) GetStatus() int32 {
	return int32(this.Asta)
}

func NewPlayer(id string, pos Vector2) *Player {
	p := new(Player)
	p.SetID(id)
	p.SetHp(100)
	p.SetPos(pos)
	p.SetSpeed(4.0)
	return p
}

// 移动
type PlayerMove struct {
	username string
	velocity *Vector2
}

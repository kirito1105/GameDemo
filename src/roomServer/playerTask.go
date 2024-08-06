package roomServer

import (
	"encoding/base64"
	"fmt"
	"google.golang.org/protobuf/proto"
	"myGameDemo/myMsg"
	"myGameDemo/mynet"
	"myGameDemo/world"
	"net"
)

type PlayerTask struct {
	username  string
	isChecked bool
	tcpTask   *mynet.TCPTask
	location  world.Point
	inRoom    *Room
}

func (p *PlayerTask) ParseMsg(data []byte) bool {
	//fmt.Println("receive msg:", data)
	//fmt.Println("receive msg:", string(data))
	var msg myMsg.Msg
	_ = proto.Unmarshal(data, &msg)
	//验证身份
	if !p.isChecked && msg.MsgType == myMsg.Type_Authentication {
		decoded, _ := base64.StdEncoding.DecodeString(msg.GetToken())

		a := CheckToken(decoded, msg.GetUsername(), msg.GetAddr(), msg.GetRoomId())
		if a {
			p.username = msg.GetUsername()
			p.isChecked = true
			p.inRoom.playerWithName[p.username] = p
			p.location = p.inRoom.world0.GetSpawn()
			fmt.Println(msg.GetUsername(), "验证成功")
		}

		for p.inRoom.world0 == nil {
		}

		spawn := p.inRoom.GetWorld().GetSpawn()
		x, y := spawn.ToUnity()
		playInfo := &myMsg.Msg{
			MsgType: myMsg.Type_PlayerInfo,
			X:       float32(x) / 100,
			Y:       float32(y) / 100,
		}
		buf, _ := proto.Marshal(playInfo)
		p.tcpTask.SendMsg(AddHeader(buf))

		for i := p.location.BlockX - 3; i <= p.location.BlockX+3; i++ {
			for j := p.location.BlockY - 3; j <= p.location.BlockY+3; j++ {
				p.SendMsgBlock(i, j)
			}
		}

	}
	//未验证的其他信息不做处理
	if !p.isChecked {
		return false
	}

	//TODO implement
	return true
}

func (p *PlayerTask) SendMsgBlock(x int, y int) { //x,y为Block坐标
	if x < 0 || x >= world.Size {
		return
	}
	if y < 0 || y >= world.Size {
		return
	}
	b := p.inRoom.GetMyWorld(x, y)
	list := make([]*myMsg.Obj, 0)
	for _, i := range b.Objs {
		x, y := i.Index.ToUnity()
		obj := &myMsg.Obj{
			X:       float32(x) / 100,
			Y:       float32(y) / 100,
			ObjType: i.ObjType,
		}
		list = append(list, obj)
	}
	block := myMsg.Block{
		Type: int64(b.TypeOfBlock),
		X:    float32(x*world.GRID_PER_BLOCK + 5),
		Y:    float32(y*world.GRID_PER_BLOCK + 5),
		List: list,
	}
	m := myMsg.Msg{
		MsgType: myMsg.Type_BlockInfo,
		Block:   &block,
	}
	fmt.Println(m)
	by, _ := proto.Marshal(&m)
	p.tcpTask.SendMsg(AddHeader(by))
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

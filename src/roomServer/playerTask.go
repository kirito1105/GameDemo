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

		a := CheckToken(decoded, msg.GetAuthentication().GetUsername(), msg.GetAuthentication().GetAddr(), msg.GetAuthentication().GetRoomId())
		if a {
			p.username = msg.GetAuthentication().GetUsername()
			p.isChecked = true
			p.inRoom.taskWithName[p.username] = p
			if p.inRoom.players[p.username] == nil {
				p.inRoom.players[p.username] = NewPlayer(
					p.username,
					*p.inRoom.world0.GetSpawn().ToVector(),
				)
			}
			fmt.Println(msg.GetAuthentication().GetUsername(), "验证成功")
			GetRoomController().PlayerOnline(p)
		}

		for p.inRoom.world0 == nil {
		}

		scene := &myMsg.MsgScene{
			Blocks: make([]*myMsg.Block, 0),
			Chars:  make([]*myMsg.CharInfo, 0),
		}
		v := p.inRoom.players[p.username].GetPos()
		playInfo := &myMsg.CharInfo{
			Username: p.username,
			Index: &myMsg.LocationInfo{
				X: v.x,
				Y: v.y,
			},
			IsUser: true,
		}
		scene.Chars = append(scene.Chars, playInfo)
		location := p.inRoom.players[p.username].pos.toPoint()
		for i := location.BlockX - 3; i <= location.BlockX+3; i++ {
			for j := location.BlockY - 3; j <= location.BlockY+3; j++ {
				scene.Blocks = append(scene.Blocks, p.GetMsgBlockWithIndex(i, j))
			}
		}
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
			velocity: Vector2{
				x: msg.GetMove().GetX(),
				y: msg.GetMove().GetY(),
			},
		}
		p.inRoom.chan_PlayerMove <- move
	}
	//TODO implement
	return true
}

func (p *PlayerTask) GetMsgBlockWithIndex(x int, y int) *myMsg.Block { //x,y为Block坐标
	if x < 0 || x >= Size {
		return nil
	}
	if y < 0 || y >= Size {
		return nil
	}
	b := p.inRoom.GetMyWorld(x, y)
	list := make([]*myMsg.Obj, 0)
	for _, i := range b.Objs {
		x, y := i.Index.ToUnity()
		obj := &myMsg.Obj{
			ObjType: i.ObjType,
			Index: &myMsg.LocationInfo{
				X: float32(x) / 100,
				Y: float32(y) / 100,
			},
		}
		list = append(list, obj)
	}

	block := myMsg.Block{
		Type: b.TypeOfBlock,
		Index: &myMsg.LocationInfo{
			X: float32(x*GRID_PER_BLOCK + 5),
			Y: float32(y*GRID_PER_BLOCK + 5),
		},
		List: list,
	}

	return &block
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
}

func NewPlayer(id string, pos Vector2) *Player {
	p := new(Player)
	p.SetID(id)
	p.SetHp(100)
	p.SetPos(pos)
	p.SetSpeed(2.0)
	return p
}

// 移动
type PlayerMove struct {
	username string
	velocity Vector2
}

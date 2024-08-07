package roomServer

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"myGameDemo/myMsg"
	"myGameDemo/mynet"
	"net"
	"time"
)

type InterForPlayer interface {
}

const (
	ft = 20 * time.Millisecond
)

type Room struct {
	RoomID       string                 //room Id
	tcpListener  mynet.TCPListener      //tco
	tasks        []*PlayerTask          //conn list
	taskWithName map[string]*PlayerTask //player map
	players      map[string]*Player
	world0       *World

	//操作
	chan_PlayerMove chan *PlayerMove

	//脏数据
	chan_Char chan *myMsg.CharInfo
}

func NewRoom() *Room {
	t := &Room{
		tcpListener:  mynet.TCPListener{},
		taskWithName: make(map[string]*PlayerTask),
		players:      make(map[string]*Player),
		world0:       nil,

		chan_PlayerMove: make(chan *PlayerMove, 10),

		chan_Char: make(chan *myMsg.CharInfo, 10),
	}
	go func() {
		a := NewWorld()
		t.world0 = a
	}()
	return t
}

func (r *Room) GetTCPAddr() net.Addr {
	return r.tcpListener.Addr()
}

func (r *Room) GetWorld() *World {
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
			r.tasks = append(r.tasks, p)
			p.Start()
		}
	}()
	go func() {
		ticker := time.Tick(ft)
		for {
			select {
			case <-ticker:
				r.Update()
			}
		}
	}()
	go func() {
		r.MainLoop()
	}()
}

func (r *Room) MainLoop() {
	for {
		select {
		case move := <-r.chan_PlayerMove:
			player := r.players[move.username]
			pos := player.GetPos()
			new_pos := pos.Add(*move.velocity.MultiplyNum(player.speed * float32(ft) / float32(time.Second)))
			fmt.Println(new_pos)
			player.SetPos(*new_pos)
			ch := &myMsg.CharInfo{
				Username: player.GetID(),
				Index: &myMsg.LocationInfo{
					X: new_pos.x,
					Y: new_pos.y,
				},
			}
			r.chan_Char <- ch
		}
	}
}

func (r *Room) Update() {
	//todo 向客户端同步信息
	//保存chan
	chan_char := r.chan_Char

	//处理角色数据
	msg := &myMsg.MsgFromService{
		Scene: NewMsgScene(),
	}
SELECT:
	for {
		select {
		case player := <-chan_char:
			msg.Scene.Chars = append(msg.Scene.Chars, player)
		default:
			break SELECT
		}
	}

	bytes, _ := proto.Marshal(msg)
	if len(bytes) == 0 {
		return
	}
	bytes = AddHeader(bytes)
	for _, player := range r.taskWithName {
		player.tcpTask.SendMsg(bytes)
	}

}

func NewMsgScene() *myMsg.MsgScene {
	return &myMsg.MsgScene{
		Blocks: make([]*myMsg.Block, 0),
		Chars:  make([]*myMsg.CharInfo, 0),
	}
}

func (r *Room) GetMyWorld(x int, y int) *Block {
	//todo
	b := r.world0.GetBlock(x, y)
	return &b
}

func (r *Room) GetOnlineList() {
	//todo
	return
}

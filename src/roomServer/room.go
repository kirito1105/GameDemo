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
	ft = 100 * time.Millisecond
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

func (this *Room) PlayerIn(p *PlayerTask) {
	this.taskWithName[p.username] = p
	_, ok := this.players[p.username]
	if !ok {
		this.players[p.username] = NewPlayer(
			p.username,
			*this.GetWorld().GetSpawn().ToVector(),
		)
	}

	this.players[p.username].Online = true

	GetRoomController().PlayerOnline(p)
}

func (this *Room) PlayerOff(p *PlayerTask) {
	//todo
}

func (this *Room) GetMsgBlockWithIndex(x int, y int) *myMsg.Block { //x,y为Block坐标
	if x < 0 || x >= Size {
		return nil
	}
	if y < 0 || y >= Size {
		return nil
	}
	b := this.GetMyWorld(x, y)
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

func (this *Room) GetInitInfo(username string) *myMsg.MsgScene {
	for this.world0 == nil {
	}
	scene := &myMsg.MsgScene{
		Blocks: make([]*myMsg.Block, 0),
		Chars:  make([]*myMsg.CharInfo, 0),
	}
	pos := this.players[username].GetPos()
	//获取所有玩家角色信息
	for _, p := range this.players {
		if !pos.CanSee(p.GetPos()) {
			continue
		}
		playInfo := &myMsg.CharInfo{
			Username: p.objId,
			Index: &myMsg.LocationInfo{
				X: p.GetPos().x,
				Y: p.GetPos().y,
			},
		}
		if p.objId == username {
			playInfo.IsUser = true
		}
		scene.Chars = append(scene.Chars, playInfo)
	}
	//获取区块信息
	location := this.players[username].pos.toPoint()
	for i := location.BlockX - 3; i <= location.BlockX+3; i++ {
		for j := location.BlockY - 3; j <= location.BlockY+3; j++ {
			scene.Blocks = append(scene.Blocks, this.GetMsgBlockWithIndex(i, j))
		}
	}
	return scene
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
				r.MainLoop()
				r.Update()
			}
		}
	}()
}

func (r *Room) MainLoop() {
STEP:
	for {
		select {
		case move := <-r.chan_PlayerMove:
			player := r.players[move.username]
			pos := player.GetPos()
			if move.velocity != nil {
				newPos := pos.Add(*move.velocity.MultiplyNum(player.speed * float32(ft) / float32(time.Second)))
				player.SetPos(*newPos)
				ch := &myMsg.CharInfo{
					Username: player.GetID(),
					Index: &myMsg.LocationInfo{
						X: newPos.x,
						Y: newPos.y,
					},
					Face: &myMsg.LocationInfo{
						X: move.velocity.x,
						Y: move.velocity.y,
					},
					AStatus: myMsg.AnimatorStatus_MOVE,
				}
				r.chan_Char <- ch
			} else {
				ch := &myMsg.CharInfo{
					Username: player.GetID(),
					Index:    nil,
					AStatus:  myMsg.AnimatorStatus_STOPMOVE,
				}
				r.chan_Char <- ch
			}
		default:
			break STEP
		}
	}
}

func (r *Room) Update() {
	//todo 向客户端同步信息
	//保存chan
	chan_char := r.chan_Char

	msg := &myMsg.MsgFromService{
		Scene: NewMsgScene(),
	}
	changed := false

	//处理角色数据
	for {
		if len(chan_char) > 0 {
			changed = true
			ch := <-chan_char
			msg.Scene.Chars = append(msg.Scene.Chars, ch)
		} else {
			break
		}
	}

	bytes, _ := proto.Marshal(msg)
	if !changed {
		return
	}
	bytes = AddHeader(bytes)
	//time.Sleep(time.Millisecond * 100)
	for _, player := range r.taskWithName {
		player.tcpTask.SendMsg(bytes)
	}
	fmt.Println(msg.Scene.Chars[0].Index)

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

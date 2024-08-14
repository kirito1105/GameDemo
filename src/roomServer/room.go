package roomServer

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"myGameDemo/myMsg"
	"myGameDemo/mynet"
	"net"
	"sync"
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
	taskWithName map[string]*PlayerTask //player map
	players      map[string]*Player
	world0       *World

	//活跃检测
	inactive *Set
	pinged   *Set

	FNO int64 //帧序号
	//操作
	chan_PlayerMove  chan *PlayerMove
	mutex_PlayerMove sync.Mutex

	//脏数据
	chan_Char chan *myMsg.CharInfo
	chan_move chan *myMsg.MoveInfo
}

func NewRoom() *Room {
	t := &Room{
		tcpListener:  mynet.TCPListener{},
		taskWithName: make(map[string]*PlayerTask),
		players:      make(map[string]*Player),
		world0:       nil,

		inactive: NewSet(),
		pinged:   NewSet(),
		FNO:      1,

		chan_PlayerMove: make(chan *PlayerMove, 10),

		chan_Char: make(chan *myMsg.CharInfo, 10),
		chan_move: make(chan *myMsg.MoveInfo, 10),
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
	new_player := &myMsg.CharInfo{
		Username: p.username,
		Index: &myMsg.LocationInfo{
			X: this.players[p.username].GetPos().x,
			Y: this.players[p.username].GetPos().y,
		},
	}

	this.chan_Char <- new_player

	this.players[p.username].Online = true

	GetRoomController().PlayerOnline(p)
}

func (this *Room) PlayerOff(u string) {
	this.taskWithName[u].Close()
	delete(this.taskWithName, u)
	this.players[u].Online = false
	//todo 广播玩家离线消息
	GetRoomController().PlayerOffline(u, this.RoomID)
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
		v := i.GetPos()
		f := i.GetObjType()
		obj := &myMsg.Obj{
			Form:    int32(f.form),
			Subform: int32(f.subForm),
			Index: &myMsg.LocationInfo{
				X: v.x,
				Y: v.y,
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
			Username: p.username,
			Index: &myMsg.LocationInfo{
				X: p.GetPos().x,
				Y: p.GetPos().y,
			},
		}
		if p.username == username {
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
			p.Start()
		}
	}()
	//TODO bufTest
	go func() {
		for {
			ele := SkillEle{
				SkillExId: 123456,
				sec:       5,
				eleId:     100,
				byStep:    SKILL_STEP_START,
				value:     100,
				timer:     time.Now().Unix() + 5,
			}
			for _, p := range r.players {
				p.bufManger.add(&ele)
			}
			time.Sleep(time.Second * 10)
		}
	}()

	//TODO END

	go func() {
		statusTicker := time.Tick(time.Second)
		ticker := time.Tick(ft)
		for {
			select {
			case <-ticker:
				//t := time.Now().UnixMicro()
				r.FNO = r.FNO + 1
				r.MoveLoop()
				r.MoveUpdate()
				r.Update()
			//fmt.Println("一帧消耗了", time.Now().UnixMicro()-t)
			case <-statusTicker:
				for _, p := range r.players {
					p.bufManger.timeTick()
				}
			}
		}
	}()
	//活跃检测
	go func() {
		ticker := time.Tick(time.Second * 30)
		for {
			select {
			case <-ticker:
				//关闭两侧未收到消息的连接
				r.pinged.Range(func(key any, value any) bool {
					u := fmt.Sprint(key)
					r.PlayerOff(u)
					return true
				})

				//向未活跃的连接发送ping

				r.inactive.Range(func(key any, value any) bool {
					u := fmt.Sprint(key)
					msg := &myMsg.MsgFromService{
						Ping: true,
					}
					bytes, _ := proto.Marshal(msg)
					r.taskWithName[u].tcpTask.SendMsg(AddHeader(bytes))
					r.pinged.Add(u)
					return true
				})
				r.inactive.Clear()
				//向一级列表加入所有其他连接
				for u, _ := range r.taskWithName {
					if !r.pinged.Exist(u) {
						r.inactive.Add(u)
					}
				}

			}
		}
	}()
}

func (r *Room) MoveLoop() {
	r.mutex_PlayerMove.Lock()
	defer r.mutex_PlayerMove.Unlock()
STEP:
	for {
		select {
		case move := <-r.chan_PlayerMove:
			player := r.players[move.username]
			pos := player.GetPos()
			if move.velocity != nil {
				newPos := pos.Add(*move.velocity.MultiplyNum(player.GetSpeed() * float32(ft) / float32(time.Second)))
				player.SetPos(*newPos)
				//ch := &myMsg.CharInfo{
				//	Username: player.GetID(),
				//	Index: &myMsg.LocationInfo{
				//		X: newPos.x,
				//		Y: newPos.y,
				//	},
				//	Face: &myMsg.LocationInfo{
				//		X: move.velocity.x,
				//		Y: move.velocity.y,
				//	},
				//	AStatus: myMsg.AnimatorStatus_MOVE,
				//}
				//r.chan_Char <- ch
				mv := &myMsg.MoveInfo{
					Username: move.username,
					Des: &myMsg.LocationInfo{
						X: newPos.x,
						Y: newPos.y,
					},
					V: &myMsg.LocationInfo{
						X: move.velocity.x,
						Y: move.velocity.y,
					},
					Speed: player.GetSpeed(),
				}
				r.chan_move <- mv

			} else {
				mv := &myMsg.MoveInfo{
					Username: move.username,
					Des: &myMsg.LocationInfo{
						X: pos.x,
						Y: pos.y,
					},
					V: &myMsg.LocationInfo{
						X: 0,
						Y: 0,
					},
					Speed: 0,
				}
				r.chan_move <- mv
			}
		default:
			break STEP
		}
	}
}

func (r *Room) MoveUpdate() {
	//保存chan
	msg := &myMsg.MsgFromService{
		Scene: NewMsgScene(),
		FNO:   r.FNO,
	}
	changed := false

	//处理角色数据
	for {
		if len(r.chan_move) > 0 {
			changed = true
			ch := <-r.chan_move
			msg.Scene.Moves = append(msg.Scene.Moves, ch)
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
}

func (this *Room) Update() {
	//保存chan
	msg := &myMsg.MsgFromService{
		Scene: NewMsgScene(),
		FNO:   this.FNO,
	}
	changed := false

	//处理角色数据
	for {
		if len(this.chan_Char) > 0 {
			changed = true
			ch := <-this.chan_Char
			msg.Scene.Chars = append(msg.Scene.Chars, ch)
		} else {
			break
		}
	}
	if !changed {
		return
	}
	bytes, _ := proto.Marshal(msg)

	bytes = AddHeader(bytes)
	//time.Sleep(time.Millisecond * 100)
	for _, player := range this.taskWithName {
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

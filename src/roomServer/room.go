package roomServer

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"myGameDemo/myMsg"
	"myGameDemo/myRand"
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
	monsters     map[int32]*Monster
	world0       *World

	//活跃检测
	inactive *Set
	pinged   *Set

	FNO int64 //帧序号
	//操作
	chan_PlayerMove  chan *PlayerMove
	mutex_PlayerMove sync.Mutex

	chan_Skill  chan *SkillMsg
	mutex_Skill sync.Mutex

	PlayerStatusUp     map[string]*myMsg.CharInfo
	mutex_PlayerStatus sync.Mutex

	//脏数据
	chan_move    chan *myMsg.MoveInfo
	chan_tree    chan *myMsg.TreeInfo
	chan_delete  chan *myMsg.DeleteInfo
	chan_monster chan *myMsg.MonsterInfo
	chan_monMove chan *myMsg.MonsterMove
}

func NewRoom() *Room {
	t := &Room{
		tcpListener:  mynet.TCPListener{},
		taskWithName: make(map[string]*PlayerTask),
		players:      make(map[string]*Player),
		monsters:     make(map[int32]*Monster),
		world0:       nil,

		inactive: NewSet(),
		pinged:   NewSet(),
		FNO:      1,

		chan_PlayerMove: make(chan *PlayerMove, 20),
		chan_Skill:      make(chan *SkillMsg, 10),
		PlayerStatusUp:  make(map[string]*myMsg.CharInfo),

		chan_move:    make(chan *myMsg.MoveInfo, 50),
		chan_tree:    make(chan *myMsg.TreeInfo, 10),
		chan_delete:  make(chan *myMsg.DeleteInfo, 10),
		chan_monster: make(chan *myMsg.MonsterInfo, 50),
		chan_monMove: make(chan *myMsg.MonsterMove, 100),
	}
	go func() {
		a := NewWorld(t)
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
		this.players[p.username].SetRoom(this)
	}
	play := this.players[p.username]
	this.players[p.username].Online = true

	this.mutex_PlayerStatus.Lock()
	defer this.mutex_PlayerStatus.Unlock()
	if val, ok := this.PlayerStatusUp[p.username]; ok {
		val.AStatus = play.GetStatus()
		val.Hp = int32(play.GetHp())
	} else {
		ch := &myMsg.CharInfo{}
		ch.AStatus = play.GetStatus()
		ch.Username = play.username
		ch.Hp = int32(play.GetHp())
		ch.Index = &myMsg.LocationInfo{
			X: play.GetPos().x,
			Y: play.GetPos().y,
		}

		this.PlayerStatusUp[play.username] = ch

	}

	GetRoomController().PlayerOnline(p)
}

func (this *Room) PlayerOff(u string) {
	if this.taskWithName[u] == nil {
		return
	}
	this.taskWithName[u].Close()
	delete(this.taskWithName, u)
	this.players[u].Online = false

	d := &myMsg.DeleteInfo{
		Form:    myMsg.Form_PLAYER,
		SubForm: myMsg.SubForm_PLAYER_01,
		Id:      this.players[u].GetID(),
		Name:    u,
	}
	this.chan_delete <- d

	GetRoomController().PlayerOffline(u)
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
			Form:    f.form,
			Subform: f.subForm,
			Index: &myMsg.LocationInfo{
				X: v.x,
				Y: v.y,
			},
			ObjId: i.GetID(),
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
		Id:   int32(x*1000 + y),
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
			Hp:      int32(p.GetHp()),
			AStatus: p.GetStatus(),
		}
		if p.username == username {
			playInfo.IsUser = true
		}
		scene.Chars = append(scene.Chars, playInfo)
	}
	//获取区块信息
	location := this.players[username].pos.toPoint()
	for i := location.BlockX - 2; i <= location.BlockX+2; i++ {
		for j := location.BlockY - 2; j <= location.BlockY+2; j++ {
			scene.Blocks = append(scene.Blocks, this.GetMsgBlockWithIndex(i, j))
		}
	}
	return scene
}

func (this *Room) GetTCPAddr() net.Addr {
	return this.tcpListener.Addr()
}

func (this *Room) GetWorld() *World {
	return this.world0
}

func (this *Room) Start() {
	err := this.tcpListener.ListenTCP()
	if err != nil {
		return
	}
	go func() {
		for {
			conn, err := this.tcpListener.AcceptTCP()
			if err != nil {
				return
			}
			p := NewPlayerTask(conn, this)
			p.Start()
		}
	}()
	//TODO bufTest
	//go func() {
	//	for {
	//		ele := SkillEle{
	//			eleId:  100,
	//			byStep: SKILL_STEP_START,
	//			value:  0,
	//			timer:  time.Now().Unix() + 5,
	//		}
	//		for _, p := range this.players {
	//			p.bufManger.add(&ele)
	//		}
	//		time.Sleep(time.Second * 10)
	//	}
	//}()

	//TODO END

	go func() {
		statusTicker := time.Tick(time.Second)
		ticker := time.Tick(ft)
		for {
			select {
			case <-ticker:
				//t := time.Now().UnixMicro()
				this.FNO = this.FNO + 1
				this.MonsterLoop()
				this.SkillLoop()
				this.MoveLoop()
				this.MoveUpdate()
				this.Update()
			//fmt.Println("一帧消耗了", time.Now().UnixMicro()-t)
			case <-statusTicker:
				for _, p := range this.players {
					if p.Online {
						p.bufManger.timeTick()
					}
				}
				this.world0.ObjManager.TimeTick()
			}
		}
	}()
	//活跃检测
	go func() {
		ticker := time.Tick(time.Second * 1)
		for {
			select {
			case <-ticker:
				//关闭两侧未收到消息的连接
				this.pinged.Range(func(key any, value any) bool {
					u := fmt.Sprint(key)

					logrus.Debug("[TCP]关闭玩家：", u, "TCP连接")
					this.PlayerOff(u)
					return true
				})

				this.pinged.Clear()
				//向未活跃的连接发送ping

				this.inactive.Range(func(key any, value any) bool {
					u := fmt.Sprint(key)
					msg := &myMsg.MsgFromService{
						Ping: true,
					}
					bytes, _ := proto.Marshal(msg)
					this.taskWithName[u].tcpTask.SendMsg(AddHeader(bytes))
					this.pinged.Add(u)
					return true
				})
				this.inactive.Clear()
				//向一级列表加入所有其他连接
				for u, _ := range this.taskWithName {
					if !this.pinged.Exist(u) {
						this.inactive.Add(u)
					}
				}

			}
		}
	}()
}

func (this *Room) MonsterLoop() {
	//生成怪物

	if this.FNO > 100 && len(this.monsters) < 1 {
		for _, p := range this.players {
			if !p.Online {
				continue
			}
			if myRand.Intn(100) < 5 {
				pig := this.world0.ObjManager.NewObj(ObjType{form: myMsg.Form_MONSTER, subForm: myMsg.SubForm_PIG})
				pig.SetRoom(this)
				v := p.GetPos()
				randv := Vector2{
					x: float32(myRand.Intn(200)) / 10,
					y: float32(myRand.Intn(200)) / 10,
				}
				new_v := v.Add(randv)
				v = this.world0.GetVectorInWorld(*new_v)
				point := v.toPoint()
				if this.GetWorld().GetBlock(point.BlockX, point.BlockY).TypeOfBlock == myMsg.BlockType_Null {
					continue
				}
				pig.SetPos(v)
				pig.SendToNine()
				logrus.Debug("[房间]生成怪物")
				logrus.Debug(p.username)
				this.monsters[pig.GetID()] = pig.(*Monster)
			}
		}
	}

	//运行状态机
	for _, m := range this.monsters {
		m.OnExecute()
	}
}

func (this *Room) SkillLoop() {
	this.mutex_Skill.Lock()
	defer this.mutex_Skill.Unlock()
SkillStep:
	for {
		select {
		case sk := <-this.chan_Skill:
			logrus.Debug(sk)
			var skill *Skill
			if sk.skillID == 0 {
				skill = this.players[sk.username].skillManager.GetSkillByID(this.players[sk.username].GetAttackID())
			} else {
				skill = this.players[sk.username].skillManager.GetSkillByID(this.players[sk.username].GetSkillID(int(sk.skillID)))
			}

			if skill == nil {
				logrus.Error("[技能]无效技能")
				continue
			}
			if sk.step == SKILL_START {
				skill.SkillActionStart(&sk.cmd, this.players[sk.username])

			} else if sk.step == SKILL_DAMAGE {
				skill.SkillActionDamage(&sk.cmd, this.players[sk.username])
			} else if sk.step == SKILL_ANIMATION {
				skill.SkillActionEnd(&sk.cmd, this.players[sk.username])
			}
		default:
			break SkillStep
		}
	}
}

func (this *Room) MoveLoop() {
	this.mutex_PlayerMove.Lock()
STEP:
	for {
		select {
		case move := <-this.chan_PlayerMove:
			player := this.players[move.username]
			pos := player.GetPos()
			oldPoint := pos.toPoint()
			if player.IsImmobilize() {
				player.RemoveStatus(ASTATUS_MOVE)
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
				this.chan_move <- mv
				continue
			}
			if move.velocity != nil {
				player.AddStatus(ASTATUS_MOVE)
				newPos := pos.Add(*move.velocity.MultiplyNum(player.GetSpeed() * float32(ft) / float32(time.Second)))
				p := newPos.toPoint()
				if this.world0.GetBlock(p.BlockX, p.BlockY).TypeOfBlock != myMsg.BlockType_Ground {
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
					this.chan_move <- mv
					continue
				}
				this.CheckMove(*oldPoint, *newPos.toPoint(), move.username)
				player.SetPos(*newPos)
				player.SetFace(Vector2{
					x: move.velocity.x,
					y: move.velocity.y,
				})
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
					Speed: player.GetSpeed() * move.velocity.magnitude(),
				}
				this.chan_move <- mv

			} else {
				player.RemoveStatus(ASTATUS_MOVE)
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
				this.chan_move <- mv
			}
		default:
			break STEP
		}
	}
	this.mutex_PlayerMove.Unlock()

}

func (this *Room) MoveUpdate() {
	//保存chan
	msg := &myMsg.MsgFromService{
		Scene: NewMsgScene(),
		FNO:   this.FNO,
	}
	changed := false

	//处理角色数据
	for {
		if len(this.chan_move) > 0 {
			changed = true
			ch := <-this.chan_move
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
	for _, player := range this.taskWithName {
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

	this.mutex_PlayerStatus.Lock()
	for _, v := range this.PlayerStatusUp {
		changed = true
		msg.Scene.Chars = append(msg.Scene.Chars, v)
	}
	this.PlayerStatusUp = make(map[string]*myMsg.CharInfo)
	this.mutex_PlayerStatus.Unlock()

	//tree
CHAN_TREE:
	for {
		select {
		case t := <-this.chan_tree:
			msg.Scene.Trees = append(msg.Scene.Trees, t)
			changed = true
		default:
			break CHAN_TREE
		}
	}

	//delete
CHAN_Delete:
	for {
		select {
		case t := <-this.chan_delete:
			msg.Scene.Deletes = append(msg.Scene.Deletes, t)
			changed = true
		default:
			break CHAN_Delete
		}
	}

	//Monster
CHAN_Monster:
	for {
		select {
		case t := <-this.chan_monster:
			msg.Scene.Monsters = append(msg.Scene.Monsters, t)
			changed = true
		default:
			break CHAN_Monster
		}
	}

	//MonsterMov
CHAN_MonsterMove:
	for {
		select {
		case t := <-this.chan_monMove:
			msg.Scene.MonsterMove = append(msg.Scene.MonsterMove, t)
			changed = true
		default:
			break CHAN_MonsterMove
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

func (this *Room) GetMyWorld(x int, y int) *Block {
	//todo
	b := this.world0.GetBlock(x, y)
	return &b
}

func (this *Room) GetOnlineList() {
	//todo
	return
}

func (this *Room) GetOnlineNum() int64 {
	return int64(len(this.taskWithName))
}

func (this *Room) CheckMove(old Point, new Point, username string) bool {
	if old.BlockX == new.BlockX && old.BlockY == new.BlockY {
		return true
	}
	scene := &myMsg.MsgScene{
		Blocks: make([]*myMsg.Block, 0),
		Chars:  make([]*myMsg.CharInfo, 0),
	}

	//获取区块信息

	for i := new.BlockX - 2; i <= new.BlockX+2; i++ {
		for j := new.BlockY - 2; j <= new.BlockY+2; j++ {
			if i >= old.BlockX && i <= old.BlockX && j >= new.BlockY && j <= new.BlockY {
				continue
			}
			scene.Blocks = append(scene.Blocks, this.GetMsgBlockWithIndex(i, j))
		}
	}
	msg := &myMsg.MsgFromService{
		FNO:   this.FNO,
		Scene: scene,
	}
	bytes, _ := proto.Marshal(msg)
	this.taskWithName[username].tcpTask.SendMsg(AddHeader(bytes))
	return true
}

type SkillMsg struct {
	username string
	cmd      StdUserAttackCMD
	skillID  int32
	step     int8
}

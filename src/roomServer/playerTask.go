package roomServer

import (
	"encoding/base64"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"myGameDemo/myMsg"
	"myGameDemo/myRand"
	"myGameDemo/mynet"
	"net"
)

type PlayerTask struct {
	username  string
	isChecked bool
	tcpTask   *mynet.TCPTask
	inRoom    *Room
	RTT       int64
	Dead      bool
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

			mm := myMsg.MsgFromService{FNO: p.inRoom.FNO}
			ss, _ := proto.Marshal(&mm)
			p.tcpTask.SendMsg(AddHeader(ss))
			p.inRoom.PlayerIn(p)
			logrus.Debug("[TCP]", msg.GetAuthentication().GetUsername(), "验证成功")

		} else {
			fmt.Println(1)
			return true
		}

		scene := p.inRoom.GetInitInfo(p.username)
		m := myMsg.MsgFromService{
			FNO:   p.inRoom.FNO,
			Scene: scene,
		}
		bytes, _ := proto.Marshal(&m)
		p.tcpTask.SendMsg(AddHeader(bytes))
	}

	//未验证的其他信息不做处理
	if !p.isChecked {
		return false
	}

	p.inRoom.inactive.Remove(p.username)
	p.inRoom.pinged.Remove(p.username)

	//移动
	if msg.GetCmd() == myMsg.Cmd_Move {
		move := &PlayerMove{
			username: p.username,
			velocity: &Vector2{
				x: msg.GetMove().GetX(),
				y: msg.GetMove().GetY(),
			},
		}
		p.inRoom.mutex_PlayerMove.Lock()
		p.inRoom.chan_PlayerMove <- move
		p.inRoom.mutex_PlayerMove.Unlock()
	}

	//停止移动
	if msg.GetCmd() == myMsg.Cmd_StopMove {
		//move := &PlayerMove{
		//	username: p.username,
		//	velocity: nil,
		//}
		//p.inRoom.mutex_PlayerMove.Lock()
		//p.inRoom.chan_PlayerMove <- move
		//p.inRoom.mutex_PlayerMove.Unlock()
		p.inRoom.mutex_PlayerStatus.Lock()
		defer p.inRoom.mutex_PlayerStatus.Unlock()
		if val, ok := p.inRoom.PlayerStatusUp[p.username]; !ok {
			p.GetPlayer().RemoveStatus(ASTATUS_MOVE)
			ch := &myMsg.CharInfo{
				Username: p.username,
				Index: &myMsg.LocationInfo{
					X: p.GetPlayer().GetPos().x,
					Y: p.GetPlayer().GetPos().y,
				},
				Face: &myMsg.LocationInfo{
					X: p.GetPlayer().GetFace().x,
					Y: p.GetPlayer().GetFace().y,
				},
				AStatus: p.GetPlayer().GetStatus(),
			}
			p.inRoom.PlayerStatusUp[p.username] = ch
		} else {
			p.GetPlayer().RemoveStatus(ASTATUS_MOVE)

			val.AStatus = p.GetPlayer().GetStatus()
			val.Face = &myMsg.LocationInfo{
				X: p.GetPlayer().GetFace().x,
				Y: p.GetPlayer().GetFace().y,
			}
			val.Index = &myMsg.LocationInfo{
				X: p.GetPlayer().GetPos().x,
				Y: p.GetPlayer().GetPos().y,
			}
		}

	}

	//技能
	//开始帧
	if msg.GetCmd() == myMsg.Cmd_Attack {

		sr := &SkillMsg{
			username: p.username,
			cmd: StdUserAttackCMD{
				direction: Vector2{
					msg.SkillRelease.Des.X,
					msg.SkillRelease.Des.Y,
				},
				location: Vector2{
					msg.SkillRelease.Pos.X,
					msg.SkillRelease.Pos.Y,
				},
			},
			skillID: msg.SkillRelease.SkillID,
			step:    SKILL_START,
		}
		go p.inRoom.SkillRelease(sr)
	}
	//伤害帧
	if msg.GetCmd() == myMsg.Cmd_Damage {
		sr := &SkillMsg{
			username: p.username,
			cmd: StdUserAttackCMD{
				direction: Vector2{
					msg.SkillRelease.Des.X,
					msg.SkillRelease.Des.Y,
				},
				location: Vector2{
					msg.SkillRelease.Pos.X,
					msg.SkillRelease.Pos.Y,
				},
			},
			skillID: msg.SkillRelease.SkillID,
			step:    SKILL_DAMAGE,
		}
		go p.inRoom.SkillRelease(sr)

	}
	//动画结束
	if msg.GetCmd() == myMsg.Cmd_Attack_EXit {
		sr := &SkillMsg{
			username: p.username,
			step:     SKILL_ANIMATION,
			skillID:  msg.SkillRelease.SkillID,
		}

		go p.inRoom.SkillRelease(sr)

	}

	//技能学习
	if msg.GetCmd() == myMsg.Cmd_Skill_learn {
		p.inRoom.players[p.username].SkillLearn(msg.SkillLearn)
	}
	//TODO implement
	return true
}

func (p *PlayerTask) GetPlayer() *Player {
	return p.inRoom.players[p.username]
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

func (this *PlayerTask) Close() {
	//todo
	this.tcpTask.Close()
}

// 玩家局内信息
type Player struct {
	ObjBase
	username     string
	Online       bool
	lastMove     Vector2
	levelManager *LevelManager
	levelUp      int
	Waiting      bool
	skillLearns  []*myMsg.SkillLearn
}

var playerId IdManager

func NewPlayer(username string, pos Vector2) *Player {
	p := new(Player)
	p.ObjType = ObjType{form: myMsg.Form_PLAYER, subForm: myMsg.SubForm_PLAYER_01}
	p.username = username
	p.SetID(playerId.getId())
	p.SetAtkBase(10)
	p.SetMaxHp(100)
	p.SetHp(100)
	p.SetPos(pos)
	p.SetSpeedBase(4.0)

	p.levelManager = NewLevelManager()

	p.BuffMaInit(p)
	p.bufManger.initOwner(p)
	p.SkillMaInit(p)
	p.skillManager.Init(p)
	p.GetSkillManager().AddSkill(1, 0)
	p.SetAttackID(1)
	p.AddTarget(SKILL_TARGET_USER)

	//p.InvincibleOn()
	return p
}

func (this *Player) SendToNine() {
	this.room.mutex_PlayerStatus.Lock()
	defer this.room.mutex_PlayerStatus.Unlock()
	if val, ok := this.GetRoom().PlayerStatusUp[this.username]; ok {
		val.AStatus = this.GetStatus()
		val.Hp = int32(this.GetHp())
	} else {
		ch := &myMsg.CharInfo{}
		ch.AStatus = this.GetStatus()
		ch.Username = this.username
		ch.Hp = int32(this.GetHp())
		ch.Index = &myMsg.LocationInfo{
			X: this.GetPos().x,
			Y: this.GetPos().y,
		}
		this.GetRoom().PlayerStatusUp[this.username] = ch

	}
}

func (this *Player) SendFaceToNine(x float32, y float32) {
	this.room.mutex_PlayerStatus.Lock()
	defer this.room.mutex_PlayerStatus.Unlock()
	if val, ok := this.GetRoom().PlayerStatusUp[this.username]; ok {
		val.Face = &myMsg.LocationInfo{
			X: x,
			Y: y,
		}
	} else {
		ch := &myMsg.CharInfo{}
		ch.AStatus = this.GetStatus()
		ch.Username = this.username
		ch.Face = &myMsg.LocationInfo{
			X: x,
			Y: y,
		}

		this.GetRoom().PlayerStatusUp[this.username] = ch

	}
}

// 移动
type PlayerMove struct {
	username string
	velocity *Vector2
}

func (this *Player) ComputeDamage(damage int) {
	if this.isInvincible() {
		return
	}
	n := (100.0 / float32(this.GetDef()+100))
	t_damage := float32(damage) * n

	this.AddHp(int(-t_damage))
	if this.isDead() {
		this.AddStatus(ASTATUS_DEAD)
		this.Immobilize()
		this.GetRoom().taskWithName[this.username].Dead = true
	} else {
		this.AddStatus(ASTATUS_INJURED)
	}

	//this.GetRoom().mutex_PlayerStatus.Lock()
	//defer this.GetRoom().mutex_PlayerStatus.Unlock()
	//if val, ok := this.GetRoom().PlayerStatusUp[this.username]; ok {
	//	val.AStatus = this.GetStatus()
	//	val.Hp = int32(this.GetHp())
	//} else {
	//	this.GetRoom().PlayerStatusUp[this.username] = &myMsg.CharInfo{
	//		Username: this.username,
	//		AStatus:  this.GetStatus(),
	//		Hp:       int32(this.GetHp()),
	//		Index: &myMsg.LocationInfo{
	//			X: this.GetPos().x,
	//			Y: this.GetPos().y,
	//		},
	//	}
	//}
	this.SendToNineNow()
	this.RemoveStatus(ASTATUS_INJURED)

}

func (this *Player) Immobilize() {
	this.immobilize = true
	this.lastMove = Vector2{0, 0}
}

func (this *Player) SendToNineNow() {
	ch := &myMsg.CharInfo{}
	ch.AStatus = this.GetStatus()
	ch.Username = this.username
	ch.Hp = int32(this.GetHp())
	ch.Index = &myMsg.LocationInfo{
		X: this.GetPos().x,
		Y: this.GetPos().y,
	}
	msg := &myMsg.MsgFromService{
		FNO:   this.GetRoom().FNO,
		Scene: NewMsgScene(),
	}
	msg.Scene.Chars = append(msg.Scene.Chars, ch)
	by, _ := proto.Marshal(msg)
	bytes := AddHeader(by)
	for _, p := range this.GetRoom().taskWithName {
		p.tcpTask.SendMsg(bytes)
	}
}

func (this *Player) AddHp(hp int) {
	if this.invincible && hp < 0 {
		return
	}
	this.hp += hp
	if this.hp < 0 {
		this.hp = 0
	}
	this.SendToNine()
}

func (this *Player) SkillLearn(learn *myMsg.SkillLearn) {
	flag := false
	for _, s := range this.skillLearns {
		if s.SkillID == learn.SkillID && s.SkillLevel == learn.SkillLevel {
			flag = true
		}
	}
	if !flag {
		return
	}
	if learn.SkillID == -1 {
		this.AddHp(10)
		this.skillLearns = nil
		this.Waiting = false
		return
	}
	this.GetSkillManager().AddSkill(learn.SkillID, learn.SkillLevel)
	if learn.SkillID < 100 {
		this.SetAttackID(learn.SkillID)
	}

	this.skillLearns = nil
	this.Waiting = false
}

func (this *Player) GetSkillLearnList() []*myMsg.SkillLearn {
	s := make([]*myMsg.SkillLearn, 0)
	if this.GetAttackID() == 1 {
		a := myRand.Intn(len(GetAttack()))
		b := myRand.Intn(len(GetAttack()) - 1)
		if b >= a {
			b++
		}
		s = append(s, &myMsg.SkillLearn{
			SkillID:    GetAttack()[a],
			SkillLevel: 0,
		})
		s = append(s, &myMsg.SkillLearn{
			SkillID:    GetAttack()[b],
			SkillLevel: 0,
		})
		s = append(s, &myMsg.SkillLearn{
			SkillID:    -1,
			SkillLevel: 0,
		})
		return s
	}
	//已经学习技能升级
	list := this.GetSkillManager().GetSkillIDAndLevel()
	newList := make([]SkillIdAndLevel, 0)
	for _, skill := range list {
		if skill.SkillId < 300 && this.GetSkillManager().GetSkillByID(skill.SkillId).GetBase().MaxLevel > skill.Level {
			newList = append(newList, skill)
		}
	}
	if len(newList) > 0 {
		a := rand.Intn(len(newList))
		s = append(s, &myMsg.SkillLearn{
			SkillID:    newList[a].SkillId,
			SkillLevel: newList[a].Level + 1,
		})
	}

	list2 := GetPassivity()
	newList2 := make([]int32, 0)
	for _, id := range list2 {
		if this.GetSkillManager().GetSkillByID(id) == nil || this.GetSkillManager().GetSkillByID(id).GetBase().MaxLevel > this.GetSkillManager().GetSkillByID(id).GetLevel() {
			newList2 = append(newList2, id)
		}
	}
	if len(newList2) > 0 {
		a := rand.Intn(len(newList2))
		var level int32
		if this.GetSkillManager().GetSkillByID(newList2[a]) == nil {
			level = 0
		} else {
			level = this.GetSkillManager().GetSkillByID(newList2[a]).GetLevel() + 1
		}
		s = append(s, &myMsg.SkillLearn{
			SkillID:    newList2[a],
			SkillLevel: level,
		})
	}

	s = append(s, &myMsg.SkillLearn{
		SkillID:    -1,
		SkillLevel: 0,
	})
	return s
}

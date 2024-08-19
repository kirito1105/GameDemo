package roomServer

import (
	"encoding/base64"
	"fmt"
	"github.com/sirupsen/logrus"
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
	logrus.Trace("[TCP]收到消息:", &msg)
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
				AStatus: myMsg.AnimatorStatus_STOPMOVE,
			}
			p.inRoom.PlayerStatusUp[p.username] = ch
		} else {
			val.AStatus = myMsg.AnimatorStatus_STOPMOVE
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
		}
		p.inRoom.mutex_Skill.Lock()
		p.inRoom.chan_Skill <- sr
		p.inRoom.mutex_Skill.Unlock()
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
	username string
	Asta     myMsg.AnimatorStatus
	Online   bool
}

func (this *Player) GetStatus() int32 {
	return int32(this.Asta)
}

func NewPlayer(username string, pos Vector2) *Player {
	p := new(Player)
	p.ObjType = ObjType{form: myMsg.Form_PLAYER, subForm: myMsg.SubForm_PLAYER_01}
	p.username = username
	p.SetHp(100)
	p.SetPos(pos)
	p.SetSpeedBase(4.0)
	p.bufManger = NewSkillStatusManager()
	p.bufManger.initOwner(p)
	return p
}

// 移动
type PlayerMove struct {
	username string
	velocity *Vector2
}

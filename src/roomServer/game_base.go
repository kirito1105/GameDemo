package roomServer

import (
	"github.com/sirupsen/logrus"
	"reflect"
)

type Skiller interface {
	GetStatusManager() *SkillStatusManager
	BuffMaInit(ObjBaseI)
	GetSkillManager() *SkillManager
	SkillMaInit(ObjBaseI)
}

type ObjBaseI interface {
	Skiller
	SendToMe()
	SendToNine()

	SendFaceToNine(x float32, y float32)

	SendToNineNow()

	GetID() int32
	SetID(id int32)

	GetObjType() ObjType
	SetObjType(objType ObjType)

	AddTarget(SkillTarget)
	RemoveTarget(SkillTarget)
	CheckTarget(SkillTarget) bool

	GetMaxHp() int
	SetMaxHp(int)

	GetHp() int
	SetHp(int)
	AddHp(int)
	isDead() bool
	ComputeDamage(int)

	isInvincible() bool
	InvincibleOn()
	InvincibleOff()

	GetPos() Vector2
	SetPos(Vector2)

	GetFace() Vector2
	SetFace(Vector2)

	GetStatus() int32
	AddStatus(int32)
	RemoveStatus(int32)

	GetSpeed() float32
	SetSpeedBase(float32)
	AddSpeed(int32)

	GetRoom() *Room
	SetRoom(*Room)

	GetAtk() int
	SetAtkBase(atk int)

	GetDef() int
	SetDef(int)

	IsImmobilize() bool
	Immobilize()
	DisImmobilize()

	GetAttackID() int32
	SetAttackID(int32)

	GetSkillID(int) int32
	SetSkillID(int, int32)

	AddAtkA(num int32)
}

type ObjBase struct {
	bufManger    *SkillStatusManager
	skillManager *SkillManager
	objId        int32       //游戏对象具有的id
	ObjType      ObjType     //对象类型
	target       SkillTarget //技能目标类型
	hp           int         //血量
	maxHp        int         //血量上限
	invincible   bool        //无敌
	pos          Vector2     //在世界中的位置
	face         Vector2     //朝向
	speedBase    float32     //移动速度,为0则无法移动
	speedTmp     int32       //加减速 百分比
	room         *Room       //所在的世界
	atkBase      int         //攻击力基础值
	atkA         int         //攻击力基础值提升
	atkB         int         //攻击力最终值提升
	atkP         int         //攻击力百分比提升
	def          int         //防御力
	Astatus      int32       //动画状态
	immobilize   bool        //定身状态
	AttackId     int32       //普攻id
	Skill1Id     int32       //技能1的id
}

func (this *ObjBase) GetStatusManager() *SkillStatusManager {
	if this.bufManger == nil {
		logrus.Panic("[角色]技能状态管理未初始化")
	}
	return this.bufManger
}

func (this *ObjBase) GetSkillManager() *SkillManager {
	if this.skillManager == nil {
		logrus.Panic("[角色]技能管理未初始化")
	}
	return this.skillManager
}

func (this *ObjBase) Init(a ObjBaseI) {
	this.BuffMaInit(a)
	this.SkillMaInit(a)
}

func (this *ObjBase) BuffMaInit(a ObjBaseI) {
	this.bufManger = NewSkillStatusManager()
	this.bufManger.initOwner(a)
}

func (this *ObjBase) SkillMaInit(a ObjBaseI) {
	this.skillManager = NewSkillManager()
	this.skillManager.Init(a)
}

func (this *ObjBase) GetID() int32 {
	return this.objId
}

func (this *ObjBase) SetID(id int32) {
	this.objId = id
}

func (this *ObjBase) GetObjType() ObjType {
	return this.ObjType
}

func (this *ObjBase) SetObjType(objType ObjType) {
	this.ObjType = objType
}

func (this *ObjBase) AddTarget(t SkillTarget) {
	this.target = this.target | t
}
func (this *ObjBase) RemoveTarget(t SkillTarget) {
	this.target = this.target & (^t)
}
func (this *ObjBase) CheckTarget(t SkillTarget) bool {
	return this.target&t != 0
}

func (this *ObjBase) GetMaxHp() int {
	return this.maxHp
}
func (this *ObjBase) SetMaxHp(set int) {
	this.maxHp = set
}

func (this *ObjBase) AddHp(hp int) {
	if this.invincible && hp < 0 {
		return
	}
	this.hp += hp
	if this.hp < 0 {
		this.hp = 0
	}
}

func (this *ObjBase) isDead() bool {
	return !(this.hp > 0)
}

func (this *ObjBase) GetHp() int {
	return this.hp
}

func (this *ObjBase) SetHp(hp int) {
	this.hp = hp
	return
}

func (this *ObjBase) GetPos() Vector2 {
	return this.pos
}

func (this *ObjBase) SetPos(pos Vector2) {
	this.pos = pos
}

func (this *ObjBase) GetFace() Vector2 {
	return this.face
}

func (this *ObjBase) SetFace(pos Vector2) {
	this.face = pos
}

func (this *ObjBase) GetSpeed() float32 {
	if this.speedBase > -1e-6 && this.speedBase < 1e-6 {
		return 0
	}
	if this.speedBase*float32(100+this.speedTmp)/100.0 < 0 {
		return 0
	}
	return this.speedBase * float32(100+this.speedTmp) / 100.0
}

func (this *ObjBase) SetSpeedBase(speed float32) {
	this.speedBase = speed
}

func (this *ObjBase) AddSpeed(t int32) {
	this.speedTmp += t
}

func (this *ObjBase) GetStatus() int32 {
	return this.Astatus
}

func (this *ObjBase) AddStatus(status int32) {
	this.Astatus = this.Astatus | status
}
func (this *ObjBase) RemoveStatus(status int32) {
	this.Astatus = this.Astatus & (^status)
}

func (this *ObjBase) GetRoom() *Room {
	return this.room
}

func (this *ObjBase) SetRoom(room *Room) {
	this.room = room
}

func (this *ObjBase) GetAtk() int {
	return int(float32(this.atkBase+this.atkA)*(1.0+float32(this.atkP)/100.0)) + this.atkB
}

func (this *ObjBase) SetAtkBase(atk int) {
	this.atkBase = atk
}

func (this *ObjBase) GetDef() int {
	return this.def
}
func (this *ObjBase) SetDef(def int) {
	this.def = def
}

func (this *ObjBase) isInvincible() bool {
	return this.invincible
}
func (this *ObjBase) InvincibleOn() {
	this.invincible = true
}
func (this *ObjBase) InvincibleOff() {
	this.invincible = false
}

func (this *ObjBase) ComputeDamage(damage int) {
	n := (100.0 / float32(this.GetDef()+100))
	t_damage := float32(damage) * n
	this.AddHp(-1 * int(t_damage))

}

func (this *ObjBase) SendToMe() {
	logrus.Error("[TCP]调用到了未被重写的函数 SendToMe")
}
func (this *ObjBase) SendToNine() {
	a := reflect.TypeOf(this)
	logrus.Error("[TCP]调用到了未被重写的函数 SendToNine", a.Name())
}

func (this *ObjBase) IsImmobilize() bool {
	return this.immobilize
}
func (this *ObjBase) Immobilize() {
	this.immobilize = true
}
func (this *ObjBase) DisImmobilize() {
	this.immobilize = false
}

func (this *ObjBase) SendFaceToNine(x float32, y float32) {
	logrus.Error("[TCP]调用到了未被重写的函数 SendFaceToNine")
}

func (this *ObjBase) GetAttackID() int32 {
	return this.AttackId
}
func (this *ObjBase) SetAttackID(id int32) {
	this.AttackId = id
}

func (this *ObjBase) GetSkillID(num int) int32 {
	switch num {
	case 1:
		return this.Skill1Id
	}

	return -1
}
func (this *ObjBase) SetSkillID(num int, id int32) {
	switch num {
	case 1:
		this.Skill1Id = id
	}
}

func (this *ObjBase) AddAtkA(num int32) {
	this.atkA += int(num)
	if this.atkA < 0 {
		this.atkA = 0
	}
}

func (this *ObjBase) SendToNineNow() {
	a := reflect.TypeOf(this)
	logrus.Error("[TCP]调用到了未被重写的函数 SendToNineNow", a.Name(), a.String())
}

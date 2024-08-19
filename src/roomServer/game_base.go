package roomServer

import "github.com/sirupsen/logrus"

type Skiller interface {
	GetStatusManager() *SkillStatusManager
	BuffMaInit()
	GetSkillManager() *SkillManager
	SkillMaInit()
}

type ObjBaseI interface {
	Skiller

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

	isInvincible() bool
	InvincibleOn()
	InvincibleOff()

	GetPos() Vector2
	SetPos(Vector2)

	GetFace() Vector2
	SetFace(Vector2)

	GetStatus() int32

	GetSpeed() float32
	SetSpeedBase(float32)
	AddSpeed(int32)

	GetRoom() *Room
	SetRoom(*Room)

	GetAtk() int
	SetAtkBase(atk int)

	GetDef() int
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

func (this *ObjBase) Init() {
	this.bufManger = NewSkillStatusManager()
	this.skillManager = NewSkillManager()
}

func (this *ObjBase) BuffMaInit() {
	this.bufManger = NewSkillStatusManager()
}

func (this *ObjBase) SkillMaInit() {
	this.skillManager = NewSkillManager()
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
	return this.hp > 0
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
	return 0
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

package roomServer

type ObjBaseI interface {
	GetID() int32
	SetID(id int32)

	GetObjType() ObjType
	SetObjType(objType ObjType)

	AddTarget(SkillTarger)
	RemoveTarget(SkillTarger)
	CheckTarget(SkillTarger) bool

	GetHp() int
	SetHp(int)

	GetPos() Vector2
	SetPos(Vector2)

	GetStatus() int32

	GetSpeed() float32
	SetSpeedBase(float32)
	AddSpeed(int32)

	GetRoom() *Room
	SetRoom(*Room)

	GetAtk() int

	GetDef() int
}

type ObjBase struct {
	bufManger *SkillStatusManager
	objId     int32       //游戏对象具有的id
	ObjType   ObjType     //对象类型
	target    SkillTarger //技能目标类型
	hp        int         //血量
	pos       Vector2     //在世界中的位置
	speedBase float32     //移动速度,为0则无法移动
	speedTmp  int32       //加减速 百分比
	room      *Room       //所在的世界
	atkBase   int         //攻击力基础值
	atkA      int         //攻击力基础值提升
	atkB      int         //攻击力最终值提升
	atkP      int         //攻击力百分比提升
	def       int         //防御力
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

func (this *ObjBase) AddTarget(t SkillTarger) {
	this.target = this.target | t
}
func (this *ObjBase) RemoveTarget(t SkillTarger) {
	this.target = this.target & (^t)
}
func (this *ObjBase) CheckTarget(t SkillTarger) bool {
	return this.target&t == 0
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

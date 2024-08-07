package roomServer

type objBaseI interface {
	GetID() string
	SetID(id string)

	GetHp() int
	SetHp(int)

	GetPos() Vector2
	SetPos(Vector2)
}

type ObjBase struct {
	objId string
	hp    int
	pos   Vector2
	speed float32
}

func (this *ObjBase) GetID() string {
	return this.objId
}

func (this *ObjBase) SetID(id string) {
	this.objId = id
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
	return this.speed
}

func (this *ObjBase) SetSpeed(speed float32) {
	this.speed = speed
}

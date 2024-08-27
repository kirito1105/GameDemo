package roomServer

import (
	"github.com/sirupsen/logrus"
	"myGameDemo/myMsg"
	"time"
)

type StateEnum int

const (
	StateMonsterDefault StateEnum = iota
	StateMonsterAttack
	StateMonsterSearch
	StateMonsterMove
	StateMonsterInjured
	StateMonsterDead
)

type StateBaseI interface {
	OnEnter()
	OnExecute()
	OnExit()
	GetState() StateEnum
	SetTransitionEvent(func(StateEnum))
}

type StateBase struct {
	state           StateEnum
	transitionEvent func(StateEnum)
	owner           ObjBaseI
}

func (s *StateBase) SetTransitionEvent(fun func(StateEnum)) {
	s.transitionEvent = fun
}
func (s *StateBase) GetState() StateEnum {
	return s.state
}

type StateMachine struct {
	stateMap     map[StateEnum]StateBaseI
	currentState StateBaseI
}

func (this *StateMachine) TransitionState(stateEnum StateEnum) {
	if this.currentState != nil {
		this.currentState.OnExit()
	}
	this.currentState = this.stateMap[stateEnum]
	this.currentState.OnEnter()
}

func (this *StateMachine) OnExecute() {
	if this.currentState != nil {
		this.currentState.OnExecute()
	}
}

func (this *StateMachine) GetState() StateBaseI {
	return this.currentState
}

type MonsterDefault struct {
	StateBase
}

func NewMonsterDefault(owner ObjBaseI) *MonsterDefault {
	this := &MonsterDefault{}
	this.state = StateMonsterDefault
	this.owner = owner
	return this
}
func (this *MonsterDefault) OnEnter() {
}
func (this *MonsterDefault) OnExecute() {
	this.transitionEvent(StateMonsterSearch)
}
func (this *MonsterDefault) OnExit() {}

type MonsterAttack struct {
	StateBase
	sec int
	cmd *StdUserAttackCMD
}

func NewMonsterAttack(owner ObjBaseI) *MonsterAttack {
	this := &MonsterAttack{}
	this.owner = owner
	this.state = StateMonsterAttack
	this.sec = 0
	return this
}
func (this *MonsterAttack) OnEnter() {
	this.owner.AddStatus(ASTATUS_ATTACK)
	this.sec = 4
}
func (this *MonsterAttack) OnExecute() {
	skill := this.owner.GetSkillManager().GetSkillByID(this.owner.GetAttackID())
	if skill == nil {
		this.transitionEvent(StateMonsterDefault)
		return
	}
	if this.sec == 4 {
		if skill.IsCD() {
			this.transitionEvent(StateMonsterDefault)
			return
		}
		this.cmd = &StdUserAttackCMD{
			location: this.owner.GetPos(),
		}
		skill.SkillActionStart(this.cmd, this.owner)
		this.sec--
		return
	}
	this.sec--
	if this.sec == 0 {
		skill.SkillActionDamage(this.cmd, this.owner)
		skill.SkillActionEnd(this.cmd, this.owner)
		this.transitionEvent(StateMonsterDefault)
	}
}
func (this *MonsterAttack) OnExit() {
	this.owner.RemoveStatus(ASTATUS_ATTACK)
	this.cmd = nil
}

type MonsterSearch struct {
	StateBase
}

func NewMonsterSearch(owner ObjBaseI) *MonsterSearch {
	this := &MonsterSearch{}
	this.owner = owner
	this.state = StateMonsterSearch
	return this
}
func (this *MonsterSearch) OnEnter() {}
func (this *MonsterSearch) OnExecute() {
	cmd := StdUserAttackCMD{
		location: this.owner.GetPos(),
	}
	list := RangePoint(SKILL_TARGET_USER, &cmd, 20, this.owner)
	if len(list) > 0 {
		pos := this.owner.GetPos()
		getPos := list[0].GetPos()
		if getPos.Add(*pos.MultiplyNum(-1)).magnitude() < 2.2 {
			this.transitionEvent(StateMonsterAttack)
			return
		}
		this.owner.(*Monster).des = list[0]
		this.transitionEvent(StateMonsterMove)
	}
}
func (this *MonsterSearch) OnExit() {}

type MonsterMove struct {
	StateBase
	lastMove Vector2
}

func NewMonsterMove(owner ObjBaseI) *MonsterMove {
	this := &MonsterMove{}
	this.owner = owner
	this.state = StateMonsterMove
	return this
}
func (this *MonsterMove) OnEnter() {
	this.owner.AddStatus(ASTATUS_MOVE)
	this.lastMove = *new(Vector2)
}

func (this *MonsterMove) OnExecute() {
	if val, ok := this.owner.(*Monster).des.(*Player); ok {
		if !val.Online {
			this.transitionEvent(StateMonsterDefault)
			return
		}
	}
	if this.owner.isInvincible() {
		return
	}

	pos := this.owner.GetPos()
	speed := this.owner.GetSpeed()
	des := this.owner.(*Monster).des.GetPos()
	v := des.Add(*pos.MultiplyNum(-1))
	if v.magnitude() < 1 {
		this.transitionEvent(StateMonsterAttack)
		return
	}

	v = v.MultiplyNum(1 / v.magnitude())
	v = v.MultiplyNum(speed)
	new_pos := pos.Add(*v.MultiplyNum((float32(ft) / float32(time.Second))))
	this.owner.SetPos(*new_pos)

	if v.Equal(&this.lastMove) {
		return
	}
	//广播移动
	msg := &myMsg.MonsterMove{
		Id:      this.owner.GetID(),
		SubForm: this.owner.GetObjType().subForm,
		Des: &myMsg.LocationInfo{
			X: new_pos.x,
			Y: new_pos.y,
		},
		V: &myMsg.LocationInfo{
			X: v.x,
			Y: v.y,
		},
		Speed:   speed,
		AStatus: this.owner.GetStatus(),
	}
	this.owner.GetRoom().chan_monMove <- msg

	//
}
func (this *MonsterMove) OnExit() {
	this.owner.RemoveStatus(ASTATUS_MOVE)

	msg := &myMsg.MonsterInfo{
		Id: this.owner.GetID(),
		Index: &myMsg.LocationInfo{
			X: this.owner.GetPos().x,
			Y: this.owner.GetPos().y,
		},
		AStatus: this.owner.GetStatus(),
		Subform: myMsg.SubForm_PIG,
	}
	this.owner.GetRoom().chan_monster <- msg
}

type MonsterInjured struct {
	StateBase
	sec int
}

func (this *MonsterInjured) OnEnter() {
	this.owner.AddStatus(ASTATUS_INJURED)
	this.owner.SendToNine()
	this.sec = 2
}
func (this *MonsterInjured) OnExecute() {
	if this.sec < 0 {
		this.transitionEvent(StateMonsterDefault)
		return
	}
	this.sec--
}
func (this *MonsterInjured) OnExit() {
	this.owner.RemoveStatus(ASTATUS_INJURED)
}

func NewMonsterInjured(owner ObjBaseI) *MonsterInjured {
	this := &MonsterInjured{}
	this.owner = owner
	this.state = StateMonsterInjured
	return this
}

type MonsterDead struct {
	StateBase
}

func NewMonsterDead(owner ObjBaseI) *MonsterDead {
	this := &MonsterDead{}
	this.owner = owner
	this.state = StateMonsterDead
	return this
}
func (this *MonsterDead) OnEnter() {
	this.owner.AddStatus(ASTATUS_DEAD)
	this.owner.SendToNine()
}
func (this *MonsterDead) OnExecute() {
	this.transitionEvent(StateMonsterDefault)
}
func (this *MonsterDead) OnExit() {
	delete(this.owner.GetRoom().monsters, this.owner.GetID())
}

type Monster struct {
	ObjBase
	StateMachine
	des ObjBaseI
}

func NewPig() *Monster {
	pig := &Monster{}
	//StateMachine初始化
	pig.stateMap = make(map[StateEnum]StateBaseI)
	pig.stateMap[StateMonsterDefault] = NewMonsterDefault(pig)
	pig.stateMap[StateMonsterAttack] = NewMonsterAttack(pig)
	pig.stateMap[StateMonsterSearch] = NewMonsterSearch(pig)
	pig.stateMap[StateMonsterMove] = NewMonsterMove(pig)
	pig.stateMap[StateMonsterInjured] = NewMonsterInjured(pig)
	pig.stateMap[StateMonsterDead] = NewMonsterDead(pig)

	for _, s := range pig.stateMap {
		s.SetTransitionEvent(pig.TransitionState)
	}
	pig.currentState = pig.stateMap[StateMonsterDefault]

	//ObjBase 初始化
	pig.Init()
	pig.SetObjType(ObjType{
		form:    myMsg.Form_MONSTER,
		subForm: myMsg.SubForm_PIG,
	})
	pig.SetHp(100)
	pig.SetAtkBase(10)
	pig.SetSpeedBase(4.2)
	pig.SetMaxHp(100)
	pig.GetStatus()

	pig.SetAttackID(101)
	pig.GetSkillManager().AddSkill(101, 0)
	return pig
}

func (this *Monster) SendToNine() {
	if this.isDead() {
		this.AddStatus(ASTATUS_DEAD)
	}
	msg := &myMsg.MonsterInfo{
		Id: this.GetID(),
		Index: &myMsg.LocationInfo{
			X: this.GetPos().x,
			Y: this.GetPos().y,
		},
		Subform: myMsg.SubForm_PIG,
		AStatus: this.GetStatus(),
	}
	this.GetRoom().chan_monster <- msg
	if this.isDead() {
		delete(this.GetRoom().monsters, this.GetID())
		this.GetStatusManager().Clear()
	}
}
func (this *Monster) ComputeDamage(damage int) {
	n := (100.0 / float32(this.GetDef()+100))
	t_damage := float32(damage) * n

	this.AddHp(int(-t_damage))
	logrus.Debug("[Monster]被命中")
	if this.isDead() {
		this.TransitionState(StateMonsterDead)
	} else {
		if this.GetState().GetState() != StateMonsterAttack {
			this.TransitionState(StateMonsterInjured)
		}

	}

}

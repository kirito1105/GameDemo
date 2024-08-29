package roomServer

import (
	"errors"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type SkillTarget int64

const (
	// 自己
	SKILL_TARGET_SELF SkillTarget = 1 << iota
	// 玩家
	SKILL_TARGET_USER
	// NPC
	SKILL_TARGET_NPC
	// 友军
	SKILL_TARGET_FRIEND
	// 怪物
	SKILL_TARGET_MONSTETR
	// 树木
	SKILL_TARGET_TREE
)
const (
	SKILL_START int8 = iota + 1
	SKILL_ANIMATION
	SKILL_DAMAGE
)

type SkillData struct {
	id        int32
	level     int32
	lastTime  int64
	SKillStep int8
}

type Skill struct {
	data  SkillData
	base  *SkillBase
	time  int64
	mutex sync.Mutex
}

func NewSkill() *Skill {
	return &Skill{}

}

func (this *Skill) GetID() int32 {
	return this.data.id
}

func (this *Skill) GetType() int32 {
	return this.base.Type
	//todo
}

func (this *Skill) SkillActionStart(cmd *StdUserAttackCMD, obj ObjBaseI) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if time.Now().Unix()-this.time > 5 {
		this.data.SKillStep = SKILL_START
		obj.DisImmobilize()
	}
	if this.data.SKillStep == SKILL_START {
		if this.IsCD() {
			logrus.Debug("[技能]", this.GetID(), "正在CD中")
			return
		}
		this.time = time.Now().Unix()
		obj.AddStatus(ASTATUS_ATTACK)
		obj.Immobilize()
		obj.RemoveStatus(ASTATUS_MOVE)
		logrus.Debug("[技能]释放技能", this.GetID())
		list := this.base.fun(this.base.Target, cmd, this.base.max, obj)
		logrus.Debug("[技能]预命中", list)
		//唯一目标技能调整朝向
		if len(list) == 1 {
			x := list[0].GetPos().x - obj.GetPos().x
			y := list[0].GetPos().y - obj.GetPos().y
			//对于obj为玩家的情况
			obj.SendFaceToNine(x, y)
		}
		//技能释放动画'
		obj.SendToNineNow()

		this.data.lastTime = time.Now().UnixMilli()
		this.data.SKillStep = SKILL_ANIMATION
	}

}

func (this *Skill) SkillActionDamage(cmd *StdUserAttackCMD, obj ObjBaseI) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.data.SKillStep == SKILL_ANIMATION {
		list := this.base.fun(this.base.Target, cmd, this.base.max, obj)
		logrus.Debug("[技能]命中", list)
		obj.RemoveStatus(ASTATUS_ATTACK)
		obj.SendToNineNow()
		skillRelease := NewSkillRelease(cmd, this, obj, list)
		skillRelease.Release()

		this.data.SKillStep = SKILL_DAMAGE
	}
}

func (this *Skill) SkillActionEnd(cmd *StdUserAttackCMD, obj ObjBaseI) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.data.SKillStep == SKILL_DAMAGE {
		this.data.SKillStep = SKILL_START
		obj.DisImmobilize()
	}
}

func (this *Skill) SkillActionClear(cmd *StdUserAttackCMD, obj ObjBaseI) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.data.SKillStep == SKILL_START {
		return
	}
	if this.data.SKillStep == SKILL_ANIMATION {
		obj.DisImmobilize()
		this.data.SKillStep = SKILL_START
		obj.RemoveStatus(ASTATUS_ATTACK)
	}
	if this.data.SKillStep == SKILL_DAMAGE {
		this.data.SKillStep = SKILL_START
		obj.DisImmobilize()
	}

}

func (this *Skill) GetCD() int64 {
	return this.base.cd
}

func (this *Skill) GetCurCd() int64 {
	n := time.Now().UnixMilli() - this.data.lastTime
	if n >= this.GetCD() {
		return 0
	}
	return this.GetCD() - n
}

func (this *Skill) SetCD() {
	this.data.lastTime = time.Now().UnixMilli()
}

func (this *Skill) IsCD() bool {
	n := time.Now().UnixMilli() - this.data.lastTime
	return !(n >= this.GetCD())
}

func (this *Skill) ClrCD() {
	this.data.lastTime = 0
}

func (this *Skill) GetLevel() int32 {
	return this.data.level
}

func (this *Skill) GetBase() SkillBase {
	return *this.base
}

type SkillList struct{}

var skList *SkillList
var oncesklist sync.Once

func GetSkillList() *SkillList {
	oncesklist.Do(func() {
		skList = &SkillList{}
	})
	return skList
}

func (this *SkillList) CreateSkill(id int32, level int32) (*Skill, error) {
	//todo
	if SkillBaseList[id] == nil {
		err := errors.New("[技能]找不到对应技能")
		return nil, err
	}
	skill := NewSkill()
	skill.base = SkillBaseList[id]
	skill.data = SkillData{
		id:        id,
		level:     level,
		lastTime:  0,
		SKillStep: SKILL_START,
	}
	return skill, nil
}

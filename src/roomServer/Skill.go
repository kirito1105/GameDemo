package roomServer

import (
	"errors"
	"sync"
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

type SkillData struct {
	id       int32
	level    int32
	lastTime int64
}

type Skill struct {
	data SkillData
	base *SkillBase
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

func (this *Skill) SkillAction(cmd StdUserAttackCMD, obj ObjBaseI) {
	list := this.base.fun(this.base.Target, cmd, this.base.max, obj)
	if list == nil {
		//有目标技能无目标
		return
	}

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
		id:       id,
		level:    level,
		lastTime: 0,
	}
	return skill, nil
}

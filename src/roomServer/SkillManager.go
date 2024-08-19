package roomServer

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	SKILL_TYPE_PASSIVITY int32 = 1 << iota
)

type SkillMap map[int32]*Skill

type SkillManager struct {
	owner    ObjBaseI
	SkillMap SkillMap
}

func NewSkillManager() *SkillManager {
	return &SkillManager{
		SkillMap: make(SkillMap),
	}
}

func (this *SkillManager) Init(owner ObjBaseI) {
	this.owner = owner
}

func (this *SkillManager) GetSkillByID(id int32) *Skill {
	if v, ok := this.SkillMap[id]; ok {
		return v
	} else {
		return nil
	}
}

func (this *SkillManager) AddSkill(id int32, level int32) error {
	skill, err := GetSkillList().CreateSkill(id, level)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	curSkill := this.GetSkillByID(id)
	//todo 对原本技能处理
	fmt.Sprint(curSkill)

	//防止被动技能叠加，删除后加入
	this.RemoveSkill(skill.GetID())
	this.runForverPassiveSkill(skill)

	this.SkillMap[id] = skill
	return nil
}

func (this *SkillManager) RemoveSkill(id int32) {
	_, ok := this.SkillMap[id]
	if !ok {
		return
	}
	//清除技能效果
	if this.owner != nil {
		this.owner.GetStatusManager().removeBySkillID(id)
	}
	delete(this.SkillMap, id)
}

func (this *SkillManager) runForverPassiveSkill(skill *Skill) {
	if skill == nil {
		return
	}
	if skill.GetType()&SKILL_TYPE_PASSIVITY != 0 {
		//todo 释放技能
	}
}

package roomServer

import (
	"github.com/sirupsen/logrus"
	"time"
)

type SkillRelease struct {
	cmd   *StdUserAttackCMD
	skill *Skill
	atk   ObjBaseI
	def   []ObjBaseI
}

func NewSkillRelease(cmd *StdUserAttackCMD, skill *Skill, atk ObjBaseI, def []ObjBaseI) *SkillRelease {
	return &SkillRelease{
		cmd:   cmd,
		skill: skill,
		atk:   atk,
		def:   def,
	}
}

func (this *SkillRelease) Release() {
	//计算直接伤害
	logrus.Debug("atk user:", this.atk.GetID())
	logrus.Debug("ATK:", this.atk.GetAtk())
	logrus.Debug("Da:", this.skill.GetBase().Damages[this.skill.GetLevel()])
	logrus.Debug("Da:", int(this.skill.GetBase().Damages[this.skill.GetLevel()]))
	damage := int(float32(this.atk.GetAtk()) * float32(this.skill.GetBase().Damages[this.skill.GetLevel()]) / 100)
	for _, d := range this.def {
		d.ComputeDamage(damage)
	}

	//施加buff
	this.doEle()

}

func (this *SkillRelease) doEle() {
	for i, id := range this.skill.GetBase().Buffs {
		var ele SkillEle
		ele.eleId = id
		ele.skillID = this.skill.GetID()
		ele.level = this.skill.GetLevel()
		ele.value = this.skill.GetBase().BuffValue[i][ele.level]
		ele.timer = time.Now().Unix() + this.skill.GetBase().BuffTime[i][ele.level]
		ele.sec = int32(this.skill.GetBase().BuffTime[i][ele.level])
		ele.byStep = SKILL_STEP_START

		ele.attackId = string(this.atk.GetID())

		for _, d := range this.def {
			this.AddAStatus(ele, d)
		}
	}
}

func (this *SkillRelease) AddAStatus(ele SkillEle, d ObjBaseI) {
	d.GetStatusManager().add(&ele)
}

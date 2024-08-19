package roomServer

type SkillRelease struct {
	cmd   *StdUserAttackCMD
	skill *Skill
	atk   ObjBaseI
	def   ObjBaseI
}

func NewSkillRelease(cmd *StdUserAttackCMD, skill *Skill, atk ObjBaseI, def ObjBaseI) *SkillRelease {
	return &SkillRelease{}
}
func (this *SkillRelease) Release() {

}

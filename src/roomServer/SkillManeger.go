package roomServer

type SkillMap map[int64]*Skill

type SkillManeger struct {
	owner    ObjBaseI
	SkillMap SkillMap
}

func NewSkillManeger() *SkillManeger {
	return &SkillManeger{
		SkillMap: make(SkillMap),
	}
}

func (this *SkillManeger) Init(owner ObjBaseI) {
	this.owner = owner
}

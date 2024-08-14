package roomServer

type SkillTarger int64

const (
	// 自己
	SKILL_TARGET_SELF SkillTarger = 1 << iota
	// 玩家
	SKILL_TARGET_USER
	// NPC
	SKILL_TARGET_NPC
	// 友军
	SKILL_TARGET_FRIEND
	// 树木
	SKILL_TARGET_TREE
)

type SkillData struct {
	id       int32
	level    int32
	lastTime int64
}

type Skill struct {
}

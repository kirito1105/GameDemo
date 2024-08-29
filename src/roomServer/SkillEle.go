package roomServer

const MAX_SKILLSTATUS = 1024

type SkillStatus int

const (
	SKILL_STATUS_TEMP      SkillStatus = iota + 1 // 临时性的
	SKILL_STATUS_TICK                             // 每秒执行状态
	SKILL_STATUS_TIME                             // 时效状态
	SKILL_STATUS_PASSIVITY                        // 永久被动

)

type SkillStep byte

const (
	SKILL_STEP_START SkillStep = iota
	SKILL_STEP_RELOAD
	SKILL_STEP_TIME
	SKILL_STEP_STOP
	SKILL_STEP_CLEAR
)

type SkillEle struct {
	eleId      int16     //buf ID
	value      int32     //技能影响的数值
	byStep     SkillStep // 状态元素类型的当前状态
	sec        int32     // 剩余秒数
	timer      int64     // 到期时间
	skillID    int32     // 技能ID
	level      int32     // 技能等级
	attackType byte      // 攻击者类型
	attackId   string    //攻击者id

}

// 空
func SkillStatus_0(obj ObjBaseI, ele *SkillEle) SkillStatus {
	if obj == nil {
		return 0
	}
	switch ele.byStep {
	case SKILL_STEP_START:
		fallthrough
	case SKILL_STEP_RELOAD:
		fallthrough
	case SKILL_STEP_TIME:
		fallthrough
	case SKILL_STEP_STOP:
		fallthrough
	case SKILL_STEP_CLEAR:

	}
	return 0
}

// 时效型加速
func SkillStatus_100(obj ObjBaseI, ele *SkillEle) SkillStatus {
	if obj == nil {
		return SKILL_STATUS_TIME
	}
	switch ele.byStep {
	case SKILL_STEP_START:
		fallthrough
	case SKILL_STEP_RELOAD:
		obj.AddSpeed(ele.value)
	case SKILL_STEP_TIME:
	case SKILL_STEP_STOP:
		fallthrough
	case SKILL_STEP_CLEAR:
		obj.AddSpeed(-ele.value)
	}
	return SKILL_STATUS_TIME
}

// 无敌
func SkillStatus_101(obj ObjBaseI, ele *SkillEle) SkillStatus {
	if obj == nil {
		return SKILL_STATUS_TIME
	}
	switch ele.byStep {
	case SKILL_STEP_START:
		fallthrough
	case SKILL_STEP_RELOAD:
		obj.InvincibleOn()
	case SKILL_STEP_TIME:
	case SKILL_STEP_STOP:
		fallthrough
	case SKILL_STEP_CLEAR:
		obj.InvincibleOff()
	}
	return SKILL_STATUS_TIME
}

// 持续掉血 凋亡（后加buff顶掉之前的凋亡）
func SkillStatus_200(obj ObjBaseI, ele *SkillEle) SkillStatus {
	if obj == nil {
		return SKILL_STATUS_TICK
	}
	switch ele.byStep {
	case SKILL_STEP_START:
		fallthrough
	case SKILL_STEP_RELOAD:
		fallthrough
	case SKILL_STEP_TIME:
		obj.AddHp(int(-ele.value))
		//logrus.Trace("[Monster]")
		if obj.isDead() {
			obj.SendToNineNow()
		}

		return SKILL_STATUS_TICK
	case SKILL_STEP_STOP:
		fallthrough
	case SKILL_STEP_CLEAR:

	}
	return 0
}

// 持续掉血 燃烧（后加buff刷新buff持续时间）
func SkillStatus_201(obj ObjBaseI, ele *SkillEle) SkillStatus {
	if obj == nil {
		return SKILL_STATUS_TICK
	}
	switch ele.byStep {
	case SKILL_STEP_START:
		fallthrough
	case SKILL_STEP_RELOAD:
		fallthrough
	case SKILL_STEP_TIME:
		obj.AddHp(int(-ele.value))
		if obj.isDead() {
			obj.SendToNineNow()
		}

		return SKILL_STATUS_TICK
	case SKILL_STEP_STOP:
		fallthrough
	case SKILL_STEP_CLEAR:

	}
	return 0
}

// ATK A up
func SkillStatus_301(obj ObjBaseI, ele *SkillEle) SkillStatus {
	if obj == nil {
		return SKILL_STATUS_TIME
	}
	switch ele.byStep {
	case SKILL_STEP_START:
		fallthrough
	case SKILL_STEP_RELOAD:
		obj.AddAtkA(ele.value)
	case SKILL_STEP_TIME:
	case SKILL_STEP_STOP:
		fallthrough
	case SKILL_STEP_CLEAR:
		obj.AddAtkA(-ele.value)
	}
	return SKILL_STATUS_TIME
}

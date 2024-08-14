package roomServer

import (
	"github.com/sirupsen/logrus"
	"time"
)

type SkillStatusMap map[int64]*SkillEle
type SkillStatusFun func(ObjBaseI, *SkillEle) SkillStatus

type SkillStatusManager struct {
	owner ObjBaseI
	//定时任务
	tickMap SkillStatusMap
	//持续
	timeMap SkillStatusMap
	//被动
	passivityMap SkillStatusMap
	//状态功能映射
	funcList []SkillStatusFun

	clearAllBuff bool
}

func NewSkillStatusManager() *SkillStatusManager {
	SSM := &SkillStatusManager{
		owner:        nil,
		tickMap:      make(SkillStatusMap),
		timeMap:      make(SkillStatusMap),
		passivityMap: make(SkillStatusMap),
		funcList:     make([]SkillStatusFun, MAX_SKILLSTATUS),
		clearAllBuff: false,
	}
	SSM.funcList[100] = SkillStatus_100
	return SSM
}

func (this *SkillStatusManager) initOwner(owner ObjBaseI) {
	this.owner = owner
}

func (this *SkillStatusManager) add(ele *SkillEle) bool {

	if this.owner == nil {
		logrus.Error("[技能]未初始化Owner")
		return false
	}
	if ele.eleId >= MAX_SKILLSTATUS {
		logrus.Error("[技能]状态ID大于最大值")
		return false
	}
	if this.funcList[ele.eleId] == nil {
		logrus.Error("[技能]找不到状态", ele.eleId)
		return false
	}

	if ele.byStep == SKILL_STEP_START {
		var oldEle *SkillEle = nil

		if _, ok := this.tickMap[ele.SkillExId]; ok {
			oldEle = this.tickMap[ele.SkillExId]
		} else if _, ok := this.timeMap[ele.SkillExId]; ok {
			oldEle = this.timeMap[ele.SkillExId]
		} else if _, ok := this.passivityMap[ele.SkillExId]; ok {
			oldEle = this.passivityMap[ele.SkillExId]
		}
		if oldEle != nil {
			oldEle.byStep = SKILL_STEP_STOP
			this.funcList[oldEle.eleId](this.owner, oldEle)

			//TODO 有特效需要通知客户端显示
			//if (oldEle.effectID > 0 && owner)
			//	owner->unShowSkillState(oldEle.skillID, oldEle.effectID);
		}
	}

	flag := true
	retCode := this.funcList[ele.eleId](this.owner, ele)
	switch retCode {
	case SKILL_STATUS_TICK:
		this.tickMap[ele.SkillExId] = ele
	case SKILL_STATUS_TIME:
		this.timeMap[ele.SkillExId] = ele
	case SKILL_STATUS_PASSIVITY:
		this.passivityMap[ele.SkillExId] = ele
	default:
		flag = false
	}
	if flag {
		// TODO 有特效需要通知客户端显示

	}
	return true
}

func (this *SkillStatusManager) timeTick() {
	tickTemp := make(SkillStatusMap)
	for k, v := range this.tickMap {
		tickTemp[k] = v
	}
	delStatus := make([]int64, 0)
	for k, v := range tickTemp {
		if v.sec <= 0 || this.clearAllBuff {
			delStatus = append(delStatus, k)
			this.SetEleStop(v)
		} else {
			v.byStep = SKILL_STEP_TIME
			this.add(v)
			v.sec--
		}

	}
	for _, v := range delStatus {
		delete(this.tickMap, v)
	}
	delStatus = make([]int64, 0)

	for k, v := range this.timeMap {
		if time.Now().Unix() > v.timer || this.clearAllBuff {
			this.SetEleStop(v)
			delStatus = append(delStatus, k)
		} else {
			v.sec--
		}
	}
	for _, v := range delStatus {
		delete(this.timeMap, v)
	}

	for _, v := range this.tickMap {
		if v.sec == 0 {
			continue
		}
		v.sec--
	}
	if this.clearAllBuff {
		this.clearAllBuff = false
	}
}

func (this *SkillStatusManager) Clear() {
	this.clearAllBuff = true
	this.timeTick()
}

func (this *SkillStatusManager) SetEleStop(ele *SkillEle) {
	ele.byStep = SKILL_STEP_STOP
	this.add(ele)

	//TODO 技能特效
}

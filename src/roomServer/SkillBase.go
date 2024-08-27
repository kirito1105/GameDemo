package roomServer

import (
	"github.com/sirupsen/logrus"
	"math"
)

var SkillBaseList []*SkillBase

const MAX_SKILL_SIZE = 1024

func init() {
	SkillBaseList = make([]*SkillBase, MAX_SKILL_SIZE)
	//平A
	SkillBaseList[1] = &SkillBase{
		MaxLevel: 0,
		Damages:  []int16{100},
		cd:       1000,
		Target:   SKILL_TARGET_TREE | SKILL_TARGET_MONSTETR,
		max:      2,
		fun:      RangePoint,
	}
	SkillBaseList[2] = &SkillBase{
		MaxLevel:   4,
		Damages:    []int16{100, 100, 110, 110, 120},
		cd:         1000,
		Target:     SKILL_TARGET_TREE | SKILL_TARGET_MONSTETR,
		Buffs:      []int16{200},
		BuffTarget: []SkillTarget{SKILL_TARGET_MONSTETR},
		BuffValue:  [][]int32{{2, 5, 8, 10, 15}},
		BuffTime:   [][]int64{{2, 2, 2, 2, 2}},
		max:        2,
		fun:        RangePoint,
	}

	//怪物平a
	SkillBaseList[101] = &SkillBase{
		MaxLevel: 4,
		Damages:  []int16{90, 100, 110, 120, 130},
		cd:       5000,
		Target:   SKILL_TARGET_USER,
		max:      2,
		fun:      RangePoint,
	}
	// 被动技能
	// ATK up
	SkillBaseList[301] = &SkillBase{
		MaxLevel:  9,
		Buffs:     []int16{301},
		BuffValue: [][]int32{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		Target:    SKILL_TARGET_USER,
		Type:      SKILL_TYPE_PASSIVITY,
	}
}

type SkillBase struct {
	MaxLevel   int16   //最大等级
	Type       int32   //技能类型
	Damages    []int16 //伤害列表
	Buffs      []int16 //buf列表
	BuffTarget []SkillTarget
	BuffValue  [][]int32   //buff数值
	BuffTime   [][]int64   //buff时间
	cd         int64       //cd
	Target     SkillTarget //技能作用对象
	Effect     int32       //特效状态id
	max        float32
	fun        func(target SkillTarget, cmd *StdUserAttackCMD, max float32, atk ObjBaseI) []ObjBaseI //技能范围
}

func RangeLine(target SkillTarget, cmd *StdUserAttackCMD, max float32, atk ObjBaseI) []ObjBaseI {
	if atk == nil {
		logrus.Error("[技能]无效的技能释放方")
		return nil
	}
	list := make([]ObjBaseI, 0)
	point := cmd.location.toPoint()

	var x, y int

	if cmd.direction.x > 0 {
		x = 1
	} else {
		x = -1
	}
	if cmd.direction.y > 0 {
		y = 1
	} else {
		y = -1
	}

	for _, o := range atk.GetRoom().GetWorld().GetBlock(point.BlockX, point.BlockY).Objs {
		if !o.CheckTarget(target) {
			continue
		}
		pos := o.GetPos()
		newPos := pos.Add(*cmd.location.MultiplyNum(-1))
		line := newPos.innerMultiply(cmd.direction)
		if line < 0 || line > max {
			continue
		}
		r := math.Acos(float64(line / newPos.magnitude()))
		sin := float32(math.Sin(r))
		dis := newPos.magnitude() * sin
		if dis > 1 {
			continue
		}
		list = append(list, o)
	}
	for _, o := range atk.GetRoom().GetWorld().GetBlock(point.BlockX+x, point.BlockY).Objs {
		if !o.CheckTarget(target) {
			continue
		}
		pos := o.GetPos()
		newPos := pos.Add(*cmd.location.MultiplyNum(-1))
		line := newPos.innerMultiply(cmd.direction)
		if line < 0 || line > max {
			continue
		}
		r := math.Acos(float64(line / newPos.magnitude()))
		sin := float32(math.Sin(r))
		dis := newPos.magnitude() * sin
		if dis > 1 {
			continue
		}
		list = append(list, o)
	}
	for _, o := range atk.GetRoom().GetWorld().GetBlock(point.BlockX, point.BlockY+y).Objs {
		if !o.CheckTarget(target) {
			continue
		}
		pos := o.GetPos()
		newPos := pos.Add(*cmd.location.MultiplyNum(-1))
		line := newPos.innerMultiply(cmd.direction)
		if line < 0 || line > max {
			continue
		}
		r := math.Acos(float64(line / newPos.magnitude()))
		sin := float32(math.Sin(r))
		dis := newPos.magnitude() * sin
		if dis > 1 {
			continue
		}
		list = append(list, o)
	}
	for _, o := range atk.GetRoom().GetWorld().GetBlock(point.BlockX+x, point.BlockY+y).Objs {
		if !o.CheckTarget(target) {
			continue
		}
		pos := o.GetPos()
		newPos := pos.Add(*cmd.location.MultiplyNum(-1))
		line := newPos.innerMultiply(cmd.direction)
		if line < 0 || line > max {
			continue
		}
		r := math.Acos(float64(line / newPos.magnitude()))
		sin := float32(math.Sin(r))
		dis := newPos.magnitude() * sin
		if dis > 1 {
			continue
		}
		list = append(list, o)
	}
	return list
}

func RangeR(target SkillTarget, cmd *StdUserAttackCMD, max float32, atk ObjBaseI) []ObjBaseI {
	point := cmd.location.toPoint()
	list := make([]ObjBaseI, 0)
	for x := point.BlockX - 1; x <= point.BlockX+1; x++ {
		for y := point.BlockY - 1; y <= point.BlockY+1; y++ {
			for _, o := range atk.GetRoom().GetWorld().GetBlock(x, y).Objs {
				if !o.CheckTarget(target) {
					continue
				}
				pos := o.GetPos()
				newPos := pos.Add(*cmd.location.MultiplyNum(-1))
				line := newPos.magnitude()
				if line < 0 || line > max {
					continue
				}
				list = append(list, o)
			}
		}
	}
	return list
}

func RangePoint(target SkillTarget, cmd *StdUserAttackCMD, max float32, atk ObjBaseI) []ObjBaseI {
	point := cmd.location.toPoint()
	var min ObjBaseI = nil
	var num = max
	if (target & SKILL_TARGET_TREE) != 0 {
		for x := point.BlockX - 1; x <= point.BlockX+1; x++ {
			for y := point.BlockY - 1; y <= point.BlockY+1; y++ {
				for _, o := range atk.GetRoom().GetWorld().GetBlock(x, y).Objs {
					if !o.CheckTarget(target) {
						continue
					}
					pos := o.GetPos()
					newPos := pos.Add(*cmd.location.MultiplyNum(-1))
					line := newPos.magnitude()
					if line < num && line > 0 {
						num = line
						min = o
					}
				}
			}
		}
	}

	if (target & SKILL_TARGET_MONSTETR) != 0 {
		for _, m := range atk.GetRoom().monsters {
			pos := m.GetPos()
			new_pos := pos.Add(*cmd.location.MultiplyNum(-1))
			line := new_pos.magnitude()
			if line < num && line > 0 {
				num = line
				min = m
			}
		}
	}
	if (target & SKILL_TARGET_USER) != 0 {
		for _, m := range atk.GetRoom().players {
			if !m.Online {
				continue
			}
			pos := m.GetPos()
			new_pos := pos.Add(*cmd.location.MultiplyNum(-1))
			line := new_pos.magnitude()
			if line < num && line > 0 {
				num = line
				min = m
			}
		}
	}
	if min == nil {
		return nil
	}
	return []ObjBaseI{min}
}

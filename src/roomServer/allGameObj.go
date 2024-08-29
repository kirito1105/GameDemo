package roomServer

import (
	"github.com/sirupsen/logrus"
	"myGameDemo/myMsg"
)

type IdManager struct {
	num int32
}

func (m *IdManager) getId() int32 {
	m.num++
	return m.num
}

type TreeObj struct {
	ObjBase
}

func NewTree() *TreeObj {
	t := &TreeObj{}
	t.target = SKILL_TARGET_TREE
	t.BuffMaInit(t)
	return t
}

type BushObj struct {
	ObjBase
}

func NewBush() *BushObj {
	b := &BushObj{}
	b.target = SKILL_TARGET_TREE
	b.BuffMaInit(b)
	return b
}

type ObjectManager struct {
	TreeList map[int32]ObjBaseI

	BushList map[int32]ObjBaseI
	AllID    IdManager
}

func NewObjManager() *ObjectManager {
	m := &ObjectManager{
		TreeList: make(map[int32]ObjBaseI),
		BushList: make(map[int32]ObjBaseI),
	}
	m.AllID.num = 1000
	return m

}

func (this *ObjectManager) NewObj(t ObjType) ObjBaseI {
	var r ObjBaseI

	switch t.form {
	case myMsg.Form_TREE:
		r = NewTree()
		r.SetID(this.AllID.getId())
		r.SetObjType(t)
		r.SetMaxHp(100)
		r.SetHp(100)
		this.TreeList[r.GetID()] = r
	case myMsg.Form_BUSH:
		r = NewBush()
		r.SetID(this.AllID.getId())
		r.SetMaxHp(100)
		r.SetHp(100)
		r.SetObjType(t)
		this.BushList[r.GetID()] = r
	case myMsg.Form_MONSTER:
		if t.subForm == myMsg.SubForm_PIG {
			r = NewPig()
			r.SetID(this.AllID.getId())
			r.SetObjType(t)
		}
		if t.subForm == myMsg.SubForm_PIG_02 {
			r = NewPig()
			r.SetID(this.AllID.getId())
			r.SetObjType(t)
			r.SetHp(60)
			r.SetDef(100)
		}
		if t.subForm == myMsg.SubForm_PIG_03 {
			r = NewPig()
			r.SetID(this.AllID.getId())
			r.SetObjType(t)
			r.SetHp(100)
			r.SetDef(50)
			r.GetSkillManager().AddSkill(102, 0)
		}

	default:
		logrus.Error("[地图]未知类型")
	}
	return r
}

func (this *ObjectManager) TimeTick() {
	for _, i := range this.TreeList {
		i.GetStatusManager().timeTick()
	}

	for _, i := range this.BushList {
		i.GetStatusManager().timeTick()
	}
}

func (this *TreeObj) ComputeDamage(damage int) {

	n := 100.0 / float32(this.GetDef()+100)
	t_damage := float32(damage) * n

	this.AddHp(int(-t_damage))
	treeinfo := &myMsg.TreeInfo{
		Id: this.GetID(),
	}
	if this.isDead() {
		treeinfo.Status = ASTATUS_DEAD
		pos := this.GetPos()
		p := pos.toPoint()
		var i int
		var obj ObjBaseI
		for i, obj = range this.GetRoom().GetWorld().GetBlock(p.BlockX, p.BlockY).Objs {
			if this.GetID() == obj.GetID() {
				break
			}
		}
		if i > len(this.GetRoom().GetWorld().blocks[p.BlockX][p.BlockY].Objs)-1 {

		} else if len(this.GetRoom().GetWorld().blocks[p.BlockX][p.BlockY].Objs) == 1 {
			this.GetRoom().GetWorld().blocks[p.BlockX][p.BlockY].Objs = nil
		} else {
			this.GetRoom().GetWorld().blocks[p.BlockX][p.BlockY].Objs = append(this.GetRoom().GetWorld().blocks[p.BlockX][p.BlockY].Objs[:i], this.GetRoom().GetWorld().blocks[p.BlockX][p.BlockY].Objs[i+1:]...)
		}

		this.GetRoom().SendEXP(50)
	} else {
		treeinfo.Status = ASTATUS_INJURED
	}

	this.GetRoom().chan_tree <- treeinfo

}

func (this *TreeObj) SendToNine() {

	treeinfo := &myMsg.TreeInfo{
		Id: this.GetID(),
	}
	if this.isDead() {
		treeinfo.Status = ASTATUS_DEAD
		this.bufManger.Clear()
		pos := this.GetPos()
		p := pos.toPoint()
		var i int
		var obj ObjBaseI
		for i, obj = range this.GetRoom().GetWorld().GetBlock(p.BlockX, p.BlockY).Objs {
			if this.GetID() == obj.GetID() {
				break
			}
		}
		if i > len(this.GetRoom().GetWorld().blocks[p.BlockX][p.BlockY].Objs)-1 {

		} else if len(this.GetRoom().GetWorld().blocks[p.BlockX][p.BlockY].Objs) == 1 {
			this.GetRoom().GetWorld().blocks[p.BlockX][p.BlockY].Objs = nil
		} else {
			this.GetRoom().GetWorld().blocks[p.BlockX][p.BlockY].Objs = append(this.GetRoom().GetWorld().blocks[p.BlockX][p.BlockY].Objs[:i], this.GetRoom().GetWorld().blocks[p.BlockX][p.BlockY].Objs[i+1:]...)
		}

		this.GetRoom().SendEXP(50)
		this.GetRoom().chan_tree <- treeinfo
	}

}
func (this *TreeObj) SendToNineNow() {
	this.SendToNine()
}

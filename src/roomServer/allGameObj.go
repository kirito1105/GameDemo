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
	return &TreeObj{}
}

type BushObj struct {
	ObjBase
}

func NewBush() *BushObj {
	return &BushObj{}
}

type MonsterObj struct {
	ObjBase
}

func NewMonsterObj() *MonsterObj {
	m := &MonsterObj{}
	m.Init()
	return m
}

type ObjectManager struct {
	TreeList map[int32]ObjBaseI

	BushList map[int32]ObjBaseI
	AllID    IdManager
}

func NewObjManager() *ObjectManager {
	return &ObjectManager{
		TreeList: make(map[int32]ObjBaseI),
		BushList: make(map[int32]ObjBaseI),
	}
}

func (this *ObjectManager) NewObj(t ObjType) ObjBaseI {
	var r ObjBaseI

	switch t.form {
	case myMsg.Form_TREE:
		r = NewTree()
		r.SetID(this.AllID.getId())
		r.SetObjType(t)
		r.BuffMaInit()
		this.TreeList[r.GetID()] = r
	case myMsg.Form_BUSH:
		r = NewBush()
		r.SetID(this.AllID.getId())
		r.SetObjType(t)
		r.BuffMaInit()
		this.BushList[r.GetID()] = r
	default:
		logrus.Error("[地图]未知类型")
	}
	return r
}

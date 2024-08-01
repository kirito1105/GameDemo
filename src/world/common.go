package world

import "sync"

type ObjNode struct {
	objType string
	name    string
	rate    int
}

var objlist *[]ObjNode
var once sync.Once

func GetList() *[]ObjNode {
	once.Do(func() {
		objlist = &[]ObjNode{
			ObjNode{"tree_01", "tree", 300},
			ObjNode{"Berry_bush_01", "BerryBush", 100},
			ObjNode{"Berry_bush_02", "TropicalBerryBush", 50},
			ObjNode{"Berry_bush_03", "JuicyBerryBush", 50},
		}
	})
	return objlist
}

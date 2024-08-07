package roomServer

import "sync"

type ObjNode struct {
	objType string
	name    string
	rate    int
}

var objlist *[]ObjNode
var once11 sync.Once

func GetList() *[]ObjNode {
	once11.Do(func() {
		objlist = &[]ObjNode{
			ObjNode{"Tree_01", "tree", 300},
			ObjNode{"Berry_bush_01", "BerryBush", 100},
			ObjNode{"Berry_bush_02", "TropicalBerryBush", 50},
			ObjNode{"Berry_bush_03", "JuicyBerryBush", 50},
		}
	})
	return objlist
}

func AddHeader(bytes []byte) []byte {
	size := len(bytes)
	buf := make([]byte, 0)
	buf = append(buf, byte(size), byte(size>>8), byte(size>>16), byte(size>>24))
	buf = append(buf, bytes...)
	return buf
}

type Vector2 struct {
	x float32
	y float32
}

func (v *Vector2) NewVector2(x, y float32) *Vector2 {
	return &Vector2{x: x, y: y}
}

func (v *Vector2) Add(vector2 Vector2) *Vector2 {

	return &Vector2{v.x + vector2.x, v.y + vector2.y}
}

func (v *Vector2) MultiplyNum(num float32) *Vector2 {
	return &Vector2{v.x * num, v.y * num}
}

func (v *Vector2) toPoint() *Point {
	x := int(v.x * 100)
	y := int(v.y * 100)
	return NewPoint(x, y)
}

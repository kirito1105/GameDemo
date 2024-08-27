package roomServer

import (
	"math"
	"myGameDemo/myMsg"
	"sync"
)

type ObjNode struct {
	objType ObjType
	name    string
	rate    int
}

type ObjType struct {
	form    myMsg.Form
	subForm myMsg.SubForm
}

var vision = 3
var objlist *[]ObjNode
var once11 sync.Once

func GetList() *[]ObjNode {
	once11.Do(func() {
		objlist = &[]ObjNode{
			{ObjType{myMsg.Form_TREE, myMsg.SubForm_Tree_01}, "tree", 100},
			//{ObjType{myMsg.Form_BUSH, myMsg.SubForm_Berry_bush_01}, "BerryBush", 20},
			//{ObjType{myMsg.Form_BUSH, myMsg.SubForm_Berry_bush_02}, "TropicalBerryBush", 10},
			//{ObjType{myMsg.Form_BUSH, myMsg.SubForm_Berry_bush_03}, "JuicyBerryBush", 10},
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

func (v *Vector2) innerMultiply(vector2 Vector2) float32 {
	return v.x*vector2.x + v.y*vector2.y
}
func (v *Vector2) magnitude() float32 {
	return float32(math.Sqrt(float64(v.x*v.x + v.y*v.y)))
}
func (v *Vector2) toPoint() *Point {
	x := int(v.x * 100)
	y := int(v.y * 100)
	return NewPoint(x, y)
}

func (v *Vector2) Equal(v2 *Vector2) bool {
	return math.Abs(float64(v.x-v2.x)) < 1e-6 && math.Abs(float64(v.y-v2.y)) < 1e-6
}

func (v *Vector2) CanSee(vector2 Vector2) bool {
	point1 := v.toPoint()
	point2 := vector2.toPoint()
	if math.Abs(float64(point1.X-point2.X)) > float64(vision) {
		return false
	}
	if math.Abs(float64(point1.Y-point2.Y)) > float64(vision) {
		return false
	}
	return true
}

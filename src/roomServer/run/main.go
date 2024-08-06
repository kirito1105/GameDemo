package main

import (
	"myGameDemo/roomServer"
)

type Info struct {
	data2 int32
}
type SimulatedSlice struct {
	array uintptr
	len   int
	cap   int
}

func main() {
	//a := Info{1}
	//Len := unsafe.Sizeof(a)
	//simSlice := &SimulatedSlice{
	//	array: uintptr(unsafe.Pointer(&a)),
	//	cap:   int(Len),
	//	len:   int(Len),
	//}
	//data := *(*[]byte)(unsafe.Pointer(simSlice))
	//fmt.Println(data)
	roomServer.GetRoomServer().Run("localhost", 2005)

}

package main

import (
	"fmt"
	_ "image/png"
	world "myGameDemo/world"
	"strconv"
	"time"
)

func main() {
	w := world.BlockCreate{}
	w.Init()
	w.ToImage("1.png")
	for i := 2; i < 6; i++ {
		w.Loop(3, 4)
		w.ToImage(strconv.Itoa(i) + ".png")
	}
	for i := 6; i < 12; i++ {
		w.Loop(0, 4)
		w.ToImage(strconv.Itoa(i) + ".png")
	}
	world1 := world.CreateWorld(w.GetWorld())
	world1.Init()

	for i := 0; i < world.Size; i++ {
		for j := 0; j < world.Size; j++ {
			if len(world1.GetBlock(i, j).Objs) > 0 {
				fmt.Println(world1.GetBlock(i, j))
			}
		}
	}
	fmt.Println(world1.GetSpawn())
	fmt.Println(world1.GetSpawn().ToUnity())
	time.Sleep(time.Minute)
}

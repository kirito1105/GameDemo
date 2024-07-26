package main

import "myGameDemo/roomServer"

func main() {
	//a := roomServer.Communication{}
	//t := a.Listen()
	//fmt.Println(t)
	//time.Sleep(time.Hour)
	roomServer.Run("localhost", 2605)
}

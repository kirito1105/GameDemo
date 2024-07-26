package rcenterServer

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"myGameDemo/myRPC"
	"myGameDemo/tokenRSA"
	"os"
	"sync"
)

func RoomServerHeart(rsInfo *myRPC.RoomServerInfo) error {
	GetRoomServerRegisterCenter().RegNewServer(rsInfo)
	return nil
}

func CreateRoom(rsInfo *myRPC.GameRoomFindInfo) (*myRPC.RoomInfo, error) {
	//TODO 创建房间

	room, err := GetRoomServerRegisterCenter().minPlayerServe().CreateRoom(context.Background(), rsInfo)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	room.Token = GetToken(rsInfo.Username, room.RoomAddr, room.RoomId)
	return room, nil
}

func GetToken(username string, addr string, roomId string) []byte {
	//todo
	byteKey, _ := os.ReadFile("rcenterServer/key.private.pem")
	var priKey rsa.PrivateKey
	err := json.Unmarshal(byteKey, &priKey)
	if err != nil {
		return nil
	}
	fmt.Println()
	str, err := tokenRSA.SignRsa(username+addr+roomId, priKey)
	if err != nil {
		return nil
	}
	return str
}

func Run() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		GetLogicRPC().server()
	}()

	wg.Wait()
}

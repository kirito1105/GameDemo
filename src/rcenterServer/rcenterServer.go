package rcenterServer

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"myGameDemo/myRPC"
	"myGameDemo/tokenRSA"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func RoomServerHeart(rsInfo *myRPC.RoomServerInfo) error {
	GetRoomServerRegisterCenter().RegNewServer(rsInfo)
	return nil
}

func CreateRoom(rsInfo *myRPC.GameRoomFindInfo) (*myRPC.RoomInfo, error) {
	//TODO 创建房间
	fmt.Println(1)
	room, err := GetRoomServerRegisterCenter().minRoomServe().CreateRoom(context.Background(), rsInfo)
	fmt.Println(rsInfo)
	fmt.Println(room)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	room.Token = GetToken(rsInfo.Username, room.RoomAddr, room.RoomId)
	return room, nil
}

func FindARoom(rsInfo *myRPC.GameRoomFindInfo) (*myRPC.RoomInfo, error) {
	room, err := GetRoomServerRegisterCenter().minRoomServe().FindARoom(context.Background(), rsInfo)
	if err != nil {
		return nil, err
	}
	room.Token = GetToken(rsInfo.Username, room.RoomAddr, room.RoomId)
	fmt.Println(room)
	return room, nil
}

func EnterRoom(info *myRPC.RoomInfoNode) (*myRPC.RoomInfo, error) {
	client := GetRoomServerRegisterCenter().FindRoomWithAddr(info.Addr)
	if client == nil {
		err := errors.New("没有对应服务器")
		return nil, err
	}
	room, err := client.EnterRoom(context.Background(), info)
	if err != nil {
		return nil, err
	}
	room.Token = GetToken(info.Username, room.RoomAddr, room.RoomId)
	return room, nil
}

func GetToken(username string, addr string, roomId string) []byte {
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

	go func() {
		GetLogicRPC().server()
	}()

	err := http.ListenAndServe("127.0.0.1:5055", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}

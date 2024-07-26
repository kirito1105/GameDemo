package roomServer

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"myGameDemo/myRPC"
	"myGameDemo/tokenRSA"
	"os"
	"strconv"
	"sync"
	"time"
)

func CheckToken(token []byte, username string, addr string, roomId string) bool {
	byteKey, _ := os.ReadFile("roomServer/key.public.pem")
	var pubKey rsa.PublicKey
	err := json.Unmarshal(byteKey, &pubKey)
	if err != nil {
		return false
	}

	flag := tokenRSA.CheckRsa(username+addr+roomId, pubKey, token)

	return flag
}

func Run(ip string, port int) {
	SetAddr(ip, port)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		GetRoomRPC().server()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			_, _ = GetMyClient().RoomServerHeart(context.Background(), &myRPC.RoomServerInfo{
				Addr:      ip + ":" + strconv.Itoa(port),
				PlayerNum: 0,
				RoomNum:   0,
			})
			time.Sleep(time.Second * 5)
		}

	}()

	wg.Wait()
}

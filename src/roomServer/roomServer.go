package roomServer

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"myGameDemo/myRPC"
	"myGameDemo/tokenRSA"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"sync"
	"time"
)

type RoomServer struct {
}

var roomServer *RoomServer
var once3 sync.Once

func GetRoomServer() *RoomServer {
	once3.Do(func() {
		roomServer = &RoomServer{}
	})
	return roomServer
}

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

func (this *RoomServer) Run(ip string, port int) {
	SetAddr(ip, port)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		GetRoomRPC().server()
	}()

	go func() {
		for {
			_, _ = GetMyClient().RoomServerHeart(context.Background(), &myRPC.RoomServerInfo{
				Addr:      ip + ":" + strconv.Itoa(port),
				PlayerNum: int64(len(GetRoomController().players)),
				RoomNum:   int64(len(GetRoomController().RoomwithId)),
			})
			time.Sleep(time.Second * 5)
		}
	}()
	http.ListenAndServe("127.0.0.1:5051", nil)
}

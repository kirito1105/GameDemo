package roomServer

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"myGameDemo/myRPC"
	"myGameDemo/tokenRSA"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"sync"
	"time"
)

func init() {
	filename := "./logs/roomServer.log"
	logfile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		fmt.Printf("[日志] 初始化失败", err)
		return
	}
	multiWriter := io.MultiWriter(os.Stdout, logfile)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceQuote:      true,                  //键值对加引号
		TimestampFormat: "2006-01-02 15:04:05", //时间格式
		FullTimestamp:   true,
	})
	logrus.SetOutput(multiWriter)
	logrus.Info("[日志]初始化完成")
}

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

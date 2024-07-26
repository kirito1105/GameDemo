package logicServer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"myGameDemo/logicServer/msg"
	"myGameDemo/logicServer/userConsole"
	"myGameDemo/myRPC"
	"net"
	"net/http"
	"strconv"
	"sync"
)

var ID string
var RPCAddr string

func register(w http.ResponseWriter, r *http.Request) {
	var auth userConsole.UserInfo
	if err := json.NewDecoder(r.Body).Decode(&auth); err != nil {
		log.Fatal(err)
		return
	}
	defer r.Body.Close()
	result, err := userConsole.GetUserConsole().Register(auth)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = msg.Send(&w, result)
	if err != nil {
		return
	}
	return
}

func login(w http.ResponseWriter, r *http.Request) {
	var auth userConsole.UserInfo
	if err := json.NewDecoder(r.Body).Decode(&auth); err != nil {
		log.Fatal(err)
		return
	}
	defer r.Body.Close()
	result, err := userConsole.GetUserConsole().Login(auth)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = msg.Send(&w, result)
	if err != nil {
		return
	}
	return

}

func getOnlineUser(w http.ResponseWriter, r *http.Request) {
	result, err := userConsole.GetUserConsole().GetOnlineUser()
	if err != nil {
		log.Fatal(err)
		return
	}
	err = msg.Send(&w, result)
	if err != nil {
		return
	}
	return
}

func getUsersList(w http.ResponseWriter, r *http.Request) {
	result, err := userConsole.GetUserConsole().GetUsersList()
	if err != nil {
		log.Fatal(err)
		return
	}
	err = msg.Send(&w, result)
	if err != nil {
		return
	}
	return
}

//func gameQueue(w http.ResponseWriter, r *http.Request) {
//	var sessionID string
//	if err := json.NewDecoder(r.Body).Decode(&sessionID); err != nil {
//		log.Fatal(err)
//		return
//	}
//	defer r.Body.Close()
//	username, err := C.GetUsername(sessionID)
//	if err != nil {
//		return
//	}
//
//}

type LogicServer struct {
	Addr string
}

func serverID(w http.ResponseWriter, r *http.Request) {
	re := msg.Res{Code: 0, Msg: ID}
	if err := json.NewEncoder(w).Encode(re); err != nil {
		log.Fatal(err)
	}

	return
}

func heart(w http.ResponseWriter, r *http.Request) {
	var session userConsole.SessionInfo
	if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
		log.Fatal(err)
		return
	}
	defer r.Body.Close()
	re, _ := userConsole.GetUserConsole().Heart(session)
	if err := json.NewEncoder(w).Encode(re); err != nil {
		log.Fatal(err)
	}

	return
}

//func matchingRoom(w http.ResponseWriter, r *http.Request) {
//	var session userConsole.SessionInfo
//	if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
//		log.Fatal(err)
//		return
//	}
//	defer r.Body.Close()
//	res, err := userConsole.GetUserConsole().SessionCheck(session)
//	if err != nil {
//		return
//	}
//	if res.Code != msg.SUCCESS {
//		msg.Send(&w, &msg.Res{Code: msg.OUTTIMESESSION, Msg: "会话过期"})
//		return
//	}
//	err = GetLogicRPC().MatchingRoom(&lrRPC.RPCUserInfo{Username: res.Msg, GameMode: 0, Addr: RPCAddr})
//	if err != nil {
//		return
//	}
//	msg.Send(&w, &msg.Res{Code: msg.SUCCESS, Msg: "已进入匹配队列"})
//}

func (N *LogicServer) Run() {

	ID = N.Addr
	portMid, _ := strconv.Atoi(ID[1:])
	RPCAddr = ":" + strconv.Itoa(portMid+111)
	//TODO Test 房间创建
	a, _ := GetRcClient().CreateRoom(context.Background(), &myRPC.GameRoomFindInfo{
		Username:   "123",
		GameMode:   myRPC.Gamemode_COOPERATION,
		MustCreate: true,
	})
	fmt.Println(a)
	conn, err := net.Dial("tcp", a.RoomAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn.Write([]byte("123"))
	//TODO End
	var wg sync.WaitGroup
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	GetLogicRPC().server(RPCAddr)
	//}()

	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/getOnlineUser", getOnlineUser)
	http.HandleFunc("/getUsersList", getUsersList)
	http.HandleFunc("/ID", serverID)
	http.HandleFunc("/heart", heart)
	//http.HandleFunc("/MatchingRoom", matchingRoom)

	if err := http.ListenAndServe(N.Addr, nil); err != nil {
		panic(err)
	}
	wg.Wait()
}

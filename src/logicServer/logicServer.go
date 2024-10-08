package logicServer

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"myGameDemo/logicServer/userConsole"
	"myGameDemo/myMsg"
	"myGameDemo/myRPC"
	"net/http"
	_ "net/http/pprof"
	"strconv"
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
	err = myMsg.Send(&w, result)
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
	err = myMsg.Send(&w, result)
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
	err = myMsg.Send(&w, result)
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
	err = myMsg.Send(&w, result)
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
	re := myMsg.Res{Code: 0, Msg: ID}
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

func createRoom(w http.ResponseWriter, r *http.Request) {
	var session userConsole.SessionInfo
	if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
		log.Fatal(err)
		return
	}
	check, err := userConsole.GetUserConsole().SessionCheck(session)
	if err != nil {
		return
	}
	if check.Code != myMsg.SUCCESS {
		myMsg.Send(&w, &myMsg.Res{Code: myMsg.OUTTIMESESSION, Msg: "会话过期"})
		return
	}
	username := check.Msg
	a, err := GetRcClient().CreateRoom(context.Background(), &myRPC.GameRoomFindInfo{
		Username:   username,
		GameMode:   myRPC.Gamemode_COOPERATION,
		MustCreate: true,
	})
	if err != nil {
		return
	}
	m := RoomToClien{
		IsFind:   a.IsFind,
		RoomAddr: a.RoomAddr,
		RoomId:   a.RoomId,
	}
	m.Token = base64.StdEncoding.EncodeToString(a.Token)
	myMsg.Send(&w, &myMsg.Res{Code: myMsg.SUCCESS, Data: m})
}

func randRoom(w http.ResponseWriter, r *http.Request) {
	var session userConsole.SessionInfo
	if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
		log.Fatal(err)
		return
	}
	check, err := userConsole.GetUserConsole().SessionCheck(session)
	if err != nil {
		return
	}
	if check.Code != myMsg.SUCCESS {
		myMsg.Send(&w, &myMsg.Res{Code: myMsg.OUTTIMESESSION, Msg: "会话过期"})
		return
	}
	username := check.Msg
	a, _ := GetRcClient().CreateRoom(context.Background(), &myRPC.GameRoomFindInfo{
		Username:   username,
		GameMode:   myRPC.Gamemode_COOPERATION,
		MustCreate: false,
	})
	m := RoomToClien{
		IsFind:   a.IsFind,
		RoomAddr: a.RoomAddr,
		RoomId:   a.RoomId,
	}
	m.Token = base64.StdEncoding.EncodeToString(a.Token)
	myMsg.Send(&w, &myMsg.Res{Code: myMsg.SUCCESS, Data: m})
}

func isLogic(w http.ResponseWriter, r *http.Request) {
	result := &myMsg.Res{
		Code: myMsg.SUCCESS,
		Msg:  "this is a logic server",
	}
	err := myMsg.Send(&w, result)
	if err != nil {
		return
	}
	return
}

func getRooms(w http.ResponseWriter, r *http.Request) {
	re, _ := GetRcClient().GetRoomList(context.Background(), &myRPC.Empty{})

	if err := json.NewEncoder(w).Encode(re); err != nil {
		log.Fatal(err)
	}

	return
}
func enterRoom(w http.ResponseWriter, r *http.Request) {
	var info myRPC.RoomInfoNode
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(info)
	room, err := GetRcClient().EnterRoom(context.Background(), &info)
	res := &myMsg.Res{}
	if err != nil {
		fmt.Println(err)
		res.Code = myMsg.NOROOM
		myMsg.Send(&w, res)
		return
	}
	res.Data = room
	res.Code = myMsg.SUCCESS
	myMsg.Send(&w, res)
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

	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/getOnlineUser", getOnlineUser)
	http.HandleFunc("/getUsersList", getUsersList)
	http.HandleFunc("/ID", serverID)
	http.HandleFunc("/heart", heart)
	http.HandleFunc("/createRoom", createRoom)
	http.HandleFunc("/isLogic", isLogic)
	http.HandleFunc("/randRoom", randRoom)
	http.HandleFunc("/getRooms", getRooms)
	http.HandleFunc("/enterRoom", enterRoom)
	//http.HandleFunc("/MatchingRoom", matchingRoom)

	if err := http.ListenAndServe(N.Addr, nil); err != nil {
		panic(err)
	}
}

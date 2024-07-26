package userConsole

type UserInfo struct {
	Username string `json:"username"`
	Pwd      string `json:"password"`
}

type SessionInfo struct {
	SessionId string `json:"sessionId"`
}

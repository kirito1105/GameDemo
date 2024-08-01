package roomServer

type MsgHead struct {
	MsgType int32
}

const (
	CHECKMSG = iota
)

package roomServer

const (
	ASTATUS_MOVE int32 = 1 << iota
	ASTATUS_ATTACK
	ASTATUS_INJURED
	ASTATUS_DEAD
)

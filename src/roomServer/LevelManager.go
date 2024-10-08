package roomServer

import "sync"

type ExpManager struct{}

func (this *ExpManager) GetExp(level int) int {
	if level > 5 && level <= 10 {
		return 200
	}
	if level > 10 {
		return 300
	}

	return 150
}

var onceEcp sync.Once
var expMa *ExpManager

func GetExpManager() *ExpManager {
	onceEcp.Do(func() {
		expMa = &ExpManager{}
	})
	return expMa
}

type LevelManager struct {
	level int
	exp   int
}

func NewLevelManager() *LevelManager {
	return &LevelManager{
		level: 0,
		exp:   100,
	}
}

func (this *LevelManager) addExp(exp int) int {
	this.exp += exp
	delt := 0
	for this.exp >= GetExpManager().GetExp(this.level) {
		delt++
		this.level++
		this.exp -= GetExpManager().GetExp(this.level)
	}
	return delt
}

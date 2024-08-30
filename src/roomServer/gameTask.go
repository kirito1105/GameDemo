package roomServer

import (
	"github.com/sirupsen/logrus"
	"myGameDemo/myMsg"
	"myGameDemo/myRand"
	"sync"
)

var gameTaskList []func(G *GameDataManager) bool
var onceTask sync.Once

func GetGameTaskList() []func(G *GameDataManager) bool {
	onceTask.Do(func() {
		gameTaskList = append(gameTaskList, Task0)
	})
	return gameTaskList
}

func Task0(G *GameDataManager) bool {
	flag := true
	if G.MonsterSkillNum > G.PigDieNum+G.Pig02DieNum+G.Pig03DieNum {
		flag = false
	}
	if G.TreeDieNumTarget > G.TreeDieNum {
		flag = false
	}
	if G.PigDieNumTarget > G.PigDieNum {
		flag = false
	}
	if G.Pig02DieNumTarget > G.Pig02DieNum {
		flag = false
	}
	if G.Pig03DieNumTarget > G.Pig03DieNum {
		flag = false
	}
	return flag
}

type GameDataManager struct {
	taskNum     int
	TreeDieNum  int
	PigDieNum   int
	Pig02DieNum int
	Pig03DieNum int

	MonsterSkillNum   int
	TreeDieNumTarget  int
	PigDieNumTarget   int
	Pig02DieNumTarget int
	Pig03DieNumTarget int
	room              *Room
}

func NewGameDataManager() *GameDataManager {
	t := &GameDataManager{}
	t.taskNum = myRand.Intn(len(GetGameTaskList()))
	t.TreeDieNumTarget = myRand.Intn(20)
	if t.TreeDieNumTarget >= 10 {
		t.TreeDieNumTarget = 0
	}
	t.PigDieNumTarget = myRand.Intn(40)
	if t.PigDieNum >= 40 {
		t.PigDieNumTarget = 0
	}

	t.Pig02DieNumTarget = myRand.Intn(20)
	if t.Pig02DieNumTarget >= 4 {
		t.Pig02DieNumTarget = 0
	}

	if t.TreeDieNumTarget+t.Pig02DieNumTarget+t.PigDieNumTarget == 0 {
		t.MonsterSkillNum = 10 + myRand.Intn(40)
	}
	return t
}

func (this *GameDataManager) InitMe(r *Room) {
	this.room = r
}

func (this *GameDataManager) GetInitInfo() *myMsg.TaskTargetInfo {
	info := &myMsg.TaskTargetInfo{
		TaskNum:  int64(this.taskNum),
		Monster:  int64(this.MonsterSkillNum),
		TreeNum:  int64(this.TreeDieNumTarget),
		Pig02Num: int64(this.Pig02DieNumTarget),
		PigNum:   int64(this.PigDieNumTarget),
	}
	return info
}

func (this *GameDataManager) Subscribe(bus *EventBus) {
	bus.Subscribe(Tree_DIE, this.TreeDie)
	bus.Subscribe(MONSTER_DIE, this.MonsterDie)
}

func (this *GameDataManager) TreeDie(IEvent) {
	this.TreeDieNum++

	this.room.mutex_Data.Lock()
	if this.room.data_Task == nil {
		this.room.data_Task = &myMsg.TaskInfo{
			TreeNum: int64(this.TreeDieNum),
		}
	} else {
		this.room.data_Task.TreeNum = int64(this.TreeDieNum)
	}
	this.room.mutex_Data.Unlock()

	if GetGameTaskList()[this.taskNum](this) {
		this.room.GameOver(true)
	}
}

func (this *GameDataManager) MonsterDie(event IEvent) {
	var data ObjType
	var ok bool
	if data, ok = event.Data().(ObjType); !ok {
		logrus.Error("[事件]未知事件")
	}
	switch data.subForm {
	case myMsg.SubForm_PIG:
		this.PigDieNum++
		this.room.mutex_Data.Lock()
		if this.room.data_Task == nil {
			this.room.data_Task = &myMsg.TaskInfo{
				PigNum: int64(this.PigDieNum),
			}
		} else {
			this.room.data_Task.PigNum = int64(this.PigDieNum)
		}
		this.room.mutex_Data.Unlock()
	case myMsg.SubForm_PIG_02:
		this.Pig02DieNum++
		this.room.mutex_Data.Lock()
		if this.room.data_Task == nil {
			this.room.data_Task = &myMsg.TaskInfo{
				Pig02Num: int64(this.Pig02DieNum),
			}
		} else {
			this.room.data_Task.Pig02Num = int64(this.Pig02DieNum)
		}
		this.room.mutex_Data.Unlock()
	case myMsg.SubForm_PIG_03:
		this.Pig03DieNum++
		this.room.mutex_Data.Lock()
		if this.room.data_Task == nil {
			this.room.data_Task = &myMsg.TaskInfo{
				Pig03Num: int64(this.Pig03DieNum),
			}
		} else {
			this.room.data_Task.Pig03Num = int64(this.Pig03DieNum)
		}
		this.room.mutex_Data.Unlock()

	}
	if GetGameTaskList()[this.taskNum](this) {
		this.room.GameOver(true)
	}
}

package servive

type ServiveInter interface {
	Init() bool
	MainLoop()
	Reload() bool
	Final() bool
}

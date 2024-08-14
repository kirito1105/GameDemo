package roomServer

import "sync"

type Set struct {
	set sync.Map
}

func (this *Set) Exist(key string) bool {
	_, exist := this.set.Load(key)
	return exist
}

func (this *Set) Add(str string) {
	if !this.Exist(str) {
		this.set.Store(str, struct{}{})
	}
}

func (this *Set) Remove(str string) {
	if this.Exist(str) {
		this.set.Delete(str)
	}
}

func (this *Set) Clear() {
	this.set = sync.Map{}
}

func (this *Set) Range(f func(key any, value any) bool) {
	this.set.Range(f)
}

func NewSet() *Set {
	return &Set{}
}

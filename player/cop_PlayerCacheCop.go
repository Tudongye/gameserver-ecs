package player

import (
	"ecs"
)

// Cop业务接口实现
type PlayerCacheCop struct {
	ecs.BaseCopV1[PlayerCache]
}

var _ ecs.Cop = &PlayerCacheCop{}

// 创建Cop时触发
func (this *PlayerCacheCop) LoadData() error {
	// key := this.GetKey()
	// data := this.GetData()
	// this.SetData(nil)
	// log.Println("PlayerCacheCop LoadData not define")
	return nil
}

// 移除Cop时触发
func (this *PlayerCacheCop) ClearData() error {
	// log.Println("PlayerCacheCop ClearData not define")
	return nil
}

// 置脏Cop，写回时触发
func (this *PlayerCacheCop) FlushData() error {
	// log.Println("PlayerCacheCop FlushData not define")
	return nil
}


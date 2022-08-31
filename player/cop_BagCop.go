package player

import (
	"ecs"
	"log"
	"errors"
)

// Cop业务接口实现
type BagCop struct {
	ecs.BaseCopV1[BagInfo]
}

var _ ecs.Cop = &BagCop{}

// 创建Cop时触发
func (this *BagCop) LoadData() error {
	key := this.GetKey()
	if bag, ok := BagInfoDBPool[key]; ok {
		this.SetData(bag)
	} else {
		return errors.New("BagCop not find")
	}
	return nil
}

// 移除Cop时触发
func (this *BagCop) ClearData() error {
	this.SetData(nil)
	return nil
}

// 置脏Cop，写回时触发
func (this *BagCop) FlushData() error {
	return nil
}

// 消耗金币
func (this *BagCop) UseGold(num int) {
	this.GetData().Gold -= num
	log.Printf("UseGold %v , Left: %v", num, this.GetData().Gold)
}
package player

import (
	"log"
)

// Entity从内存移除时触发
func (this *PlayerEntity) ClearEntity() {
	// log.Println("PlayerEntity ClearEntity not define")
}

// 创建新的Entity缓存时触发
func (this *PlayerEntity) InitEntity() error {
	log.Println("PlayerEntity InitEntity")
	return nil
}


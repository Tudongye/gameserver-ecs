package player

import (
	"log"
)

// World.Start时触发
func (this *PlayerWorld) Prepare() error {
	log.Println("PlayerWorld Prepare not define")
	return nil
}

// ECS框架日志接口
func (this *PlayerWorld) PrintECSLog(s string) {
	log.Println(s)
}


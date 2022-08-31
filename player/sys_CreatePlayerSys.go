package player

import (
	"ecs"
	"log"
)

// 路由匹配时触发
func (this *CreatePlayerSys) RouteMatch(payload interface{}) bool {
	//log.Println("CreatePlayerSys RouteMatch not define")
	msg, _ := payload.(Msg)
	if msg.Sys == "CreatePlayerSys" {
		return true
	}
	return false
}

// 运行Sys业务逻辑时触发
func (this *CreatePlayerSys) MainFunc(payload interface{}, entity ecs.Entity) error {
	log.Println("CreatePlayerSys MainFunc")
	msg, _ := payload.(Msg)
	WeaponInfoDBPool[ecs.KeyType(msg.Key)] = &WeaponInfo{}
	BagInfoDBPool[ecs.KeyType(msg.Key)] = &BagInfo{100}
	return nil
}


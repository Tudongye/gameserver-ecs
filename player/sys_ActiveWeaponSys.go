package player

import (
	"ecs"
	"log"
)

// 路由匹配时触发
func (this *ActiveWeaponSys) RouteMatch(payload interface{}) bool {
	// log.Println("ActiveWeaponSys RouteMatch not define")
	msg, _ := payload.(Msg)
	if msg.Sys == "ActiveWeaponSys" {
		return true
	}
	return false
}

// 运行Sys业务逻辑时触发
func (this *ActiveWeaponSys) MainFunc(payload interface{}, entity ecs.Entity) error {
	log.Println("ActiveWeaponSys MainFunc")
	msg, _ := payload.(Msg)
	weaponcop := this.GetWeaponCop(entity)
	bagcop := this.GetBagCop(entity)
	weaponcop.ActiviteWeapon(msg.Para, bagcop)
	
	return nil
}


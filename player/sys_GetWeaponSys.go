package player

import (
	"ecs"
	"log"
)

// 路由匹配时触发
func (this *GetWeaponSys) RouteMatch(payload interface{}) bool {
	//log.Println("GetWeaponSys RouteMatch not define")
	msg, _ := payload.(Msg)
	if msg.Sys == "GetWeaponSys" {
		return true
	}
	return false
}

// 运行Sys业务逻辑时触发
func (this *GetWeaponSys) MainFunc(payload interface{}, entity ecs.Entity) error {
	log.Println("GetWeaponSys MainFunc")
	weaponcop := this.GetWeaponCop(entity)
	weaponcop.GetWeaponList()
	return nil
}


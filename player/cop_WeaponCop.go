package player

import (
	"ecs"
	"log"
	"errors"
)

// Cop业务接口实现
type WeaponCop struct {
	ecs.BaseCopV1[WeaponInfo]
}

var _ ecs.Cop = &WeaponCop{}

// 创建Cop时触发
func (this *WeaponCop) LoadData() error {
	key := this.GetKey()
	if weapon, ok := WeaponInfoDBPool[key]; ok {
		this.SetData(weapon)
	} else {
		return errors.New("WeaponCop not find")
	}
	return nil
}

// 移除Cop时触发
func (this *WeaponCop) ClearData() error {
	this.SetData(nil)
	return nil
}

// 置脏Cop，写回时触发
func (this *WeaponCop) FlushData() error {
	return nil
}

func (this *WeaponCop) GetWeaponList() []int {
	log.Printf("GetWeaponList: %v", this.GetData().WeaponList)
	return this.GetData().WeaponList
}

func (this *WeaponCop) ActiviteWeapon(weaponid int, bagcop *BagCop) {
	this.GetData().WeaponList = append(this.GetData().WeaponList, weaponid)
	bagcop.UseGold(5)
	log.Printf("ActiviteWeapon: %v", weaponid)
}
package player

import (
	"ecs"
)

type PlayerEntity struct {
	ecs.BaseEntityV1
}

var _ ecs.Entity = &PlayerEntity{}

type PlayerEntity_Cop int

// Cop编号
const (
	PlayerEntity_Cop_WeaponCop PlayerEntity_Cop = 1
	PlayerEntity_Cop_BagCop PlayerEntity_Cop = 2
	PlayerEntity_Cop_PlayerCacheCop PlayerEntity_Cop = 3
)

var CacheCopList map[int]bool

// 缓存型Cop
func GetCacheCopList() map[int]bool {
	if CacheCopList == nil {
		CacheCopList = make(map[int]bool)
		CacheCopList[int(PlayerEntity_Cop_WeaponCop)] = true
		CacheCopList[int(PlayerEntity_Cop_BagCop)] = true
	}
	return CacheCopList
}

func (this *PlayerEntity) CreateCop(copid int, key ecs.KeyType) ecs.Cop {
	var cop ecs.Cop
	if copid < 0 {
		return nil
	} else if copid == int(PlayerEntity_Cop_WeaponCop) {
		cop = &WeaponCop{}
	} else if copid == int(PlayerEntity_Cop_BagCop) {
		cop = &BagCop{}
	} else if copid == int(PlayerEntity_Cop_PlayerCacheCop) {
		cop = &PlayerCacheCop{}
	} else {
		return nil
	}
	cop.ECS_C_Construct(cop, key)
	return cop
}


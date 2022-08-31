package player

import (
	"ecs"
)


type CreatePlayerSys struct {
	ecs.BaseSysV1
}

var _ ecs.Sys = &CreatePlayerSys{}

func (this *CreatePlayerSys) Construct() {
	this.ECS_S_Construct(this)
	this.ECS_S_RegisterCop(int(PlayerEntity_Cop_PlayerCacheCop))
}

func (this *CreatePlayerSys) GetPlayerCacheCop(e ecs.Entity) *PlayerCacheCop {
	return ecs.SysCopHelper[*PlayerCacheCop](this, e, int(PlayerEntity_Cop_PlayerCacheCop))
}



type ActiveWeaponSys struct {
	ecs.BaseSysV1
}

var _ ecs.Sys = &ActiveWeaponSys{}

func (this *ActiveWeaponSys) Construct() {
	this.ECS_S_Construct(this)
	this.ECS_S_RegisterCop(int(PlayerEntity_Cop_WeaponCop))
	this.ECS_S_RegisterCop(int(PlayerEntity_Cop_BagCop))
	this.ECS_S_RegisterCop(int(PlayerEntity_Cop_PlayerCacheCop))
}

func (this *ActiveWeaponSys) GetWeaponCop(e ecs.Entity) *WeaponCop {
	return ecs.SysCopHelper[*WeaponCop](this, e, int(PlayerEntity_Cop_WeaponCop))
}

func (this *ActiveWeaponSys) GetBagCop(e ecs.Entity) *BagCop {
	return ecs.SysCopHelper[*BagCop](this, e, int(PlayerEntity_Cop_BagCop))
}

func (this *ActiveWeaponSys) GetPlayerCacheCop(e ecs.Entity) *PlayerCacheCop {
	return ecs.SysCopHelper[*PlayerCacheCop](this, e, int(PlayerEntity_Cop_PlayerCacheCop))
}



type GetWeaponSys struct {
	ecs.BaseSysV1
}

var _ ecs.Sys = &GetWeaponSys{}

func (this *GetWeaponSys) Construct() {
	this.ECS_S_Construct(this)
	this.ECS_S_RegisterCop(int(PlayerEntity_Cop_WeaponCop))
	this.ECS_S_RegisterCop(int(PlayerEntity_Cop_PlayerCacheCop))
}

func (this *GetWeaponSys) GetWeaponCop(e ecs.Entity) *WeaponCop {
	return ecs.SysCopHelper[*WeaponCop](this, e, int(PlayerEntity_Cop_WeaponCop))
}

func (this *GetWeaponSys) GetPlayerCacheCop(e ecs.Entity) *PlayerCacheCop {
	return ecs.SysCopHelper[*PlayerCacheCop](this, e, int(PlayerEntity_Cop_PlayerCacheCop))
}




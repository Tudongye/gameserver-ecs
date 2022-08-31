package player

import (
	"ecs"
)

type PlayerWorld struct {
	ecs.BaseWorldV1
}

var _ ecs.World = &PlayerWorld{}

func (this *PlayerWorld) CreateEntity(key ecs.KeyType) ecs.Entity {
	entity := &PlayerEntity{}
	entity.ECS_E_Construct(entity, key)
	entity.InitEntity()
	return entity
}

func (this *PlayerWorld) Construct() {
	this.ECS_W_Construct(this, GetCacheCopList())
	CreatePlayerSys := &CreatePlayerSys{}
	CreatePlayerSys.Construct()
	this.ECS_W_RegisterSys(CreatePlayerSys)
	ActiveWeaponSys := &ActiveWeaponSys{}
	ActiveWeaponSys.Construct()
	this.ECS_W_RegisterSys(ActiveWeaponSys)
	GetWeaponSys := &GetWeaponSys{}
	GetWeaponSys.Construct()
	this.ECS_W_RegisterSys(GetWeaponSys)
}

func (this *PlayerWorld) Start(done chan bool) {
	this.ECS_W_Start(done)
}

func (this *PlayerWorld) PushMsg(key ecs.KeyType, msg interface{}) error {
	return this.ECS_W_Push(key, msg)
}

func CreatePlayerWorld() *PlayerWorld {
	world := &PlayerWorld{}
	world.Construct()
	world.Prepare()
	return world
}


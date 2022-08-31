package player

import (
	"ecs"
)

// 假装是DB
var WeaponInfoDBPool map[ecs.KeyType]*WeaponInfo = make(map[ecs.KeyType]*WeaponInfo)
var BagInfoDBPool map[ecs.KeyType]*BagInfo = make(map[ecs.KeyType]*BagInfo)

type WeaponInfo struct {
	WeaponList []int
}

type BagInfo struct {
	Gold int
}

type PlayerCache struct {

}

type Msg struct {
	Sys 	string 	
	Key 	int
	Para    int
	Desc    string
}
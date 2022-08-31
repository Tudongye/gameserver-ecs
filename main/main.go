package main

import (
	"player"
	"time"
	"log"
)

func main() {
	world := player.CreatePlayerWorld()
	done := make(chan bool)
	world.Start(done)
	msglist := make([]player.Msg, 0)
	msglist = append(msglist, player.Msg{"ActiveWeaponSys", 1, 1001,"给玩家1激活武器1001"})	
	msglist = append(msglist, player.Msg{"CreatePlayerSys", 1, 0,"创建玩家1"})	
	msglist = append(msglist, player.Msg{"ActiveWeaponSys", 1, 1001,"给玩家1激活武器1001"})	
	msglist = append(msglist, player.Msg{"GetWeaponSys", 1, 0,"获取玩家1得武器列表"})	
	for _, msg := range msglist {
		world.PushMsg(1, msg)
		log.Printf("Call: %v", msg)
		time.Sleep(time.Duration(1)*time.Second) // 
	}

	<-done
}

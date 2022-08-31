module main

go 1.18

replace player => ../../frame/player

replace ecs => ../../frame/ecs

require (
	ecs v0.0.0-00010101000000-000000000000 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	player v0.0.0-00010101000000-000000000000 // indirect

)

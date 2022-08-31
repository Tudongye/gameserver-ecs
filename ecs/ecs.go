// @Title  ecs.go
// @Description  ecs模式库
// @Author  panhuili 20220826
// @Update  panhuili 20220826

package ecs

import (
	"log"
)

// EcsCfg 配置
type EcsCfg struct {
	// World
	World_MaxMsgQueueLen   int // 消息队列长度
	World_MaxMsgNumPerLoop int // 单帧最大处理消息数
	World_MaxWorkQueue     int // 全局最大工作数
	World_WorkTimeOut      int // 单个工作超时时间
	World_WorkerPoolSize   int // 工作者池

	World_StatPrintInter int           // 打印World统计信息间隔
	World_LogLevel       ECS_LOG_LEVEL // 日志等级

	// Entity
	Entity_MaxWorkerNum int // Entity 单个Entity最大工作数
	Entity_ClearTimeOut int // Entity 超时未激活移除时间

	// Cop
	Cop_DirtyTimeOut int // Cop 置脏强制写回时间
	Cop_ClearTimeOut int // Cop 超时未激活移除时间
}

var cfg *EcsCfg

// @title    InitEcsCfg
// @description   初始化ECS配置
// @auth      panhuili             20220826
// @param     c        *EcsCfg         "业务自定义配置"
func InitEcsCfg(c *EcsCfg) {
	log.Println("InitEcsCfg")
	cfg = c
}

// @title    GetDefaultECSCfg
// @description   生成默认ECS配置
// @auth      panhuili             20220826
// @return     c        *EcsCfg         "默认ECS配置"
func GetDefaultECSCfg() *EcsCfg {
	return &EcsCfg{
		World_MaxMsgQueueLen:   100,
		World_MaxMsgNumPerLoop: 50,
		Entity_ClearTimeOut:    10,
		World_MaxWorkQueue:     100,
		World_WorkTimeOut:      5,
		World_WorkerPoolSize:   10,
		World_StatPrintInter:   60,
		World_LogLevel:         ECS_LOG_ERROR,
		Entity_MaxWorkerNum:    10,
		Cop_DirtyTimeOut:       20,
		Cop_ClearTimeOut:       45,
	}
}

// @title    GetEcsCfg
// @description   获取ECS配置
// @auth      panhuili             20220826
// @return     c        *EcsCfg         "ECS配置"
func GetEcsCfg() *EcsCfg {
	if cfg == nil {
		cfg = GetDefaultECSCfg()
	}
	return cfg
}

// ECS日志等级
type ECS_LOG_LEVEL int

const (
	ECS_LOG_DEBUG ECS_LOG_LEVEL = 1
	ECS_LOG_WARN  ECS_LOG_LEVEL = 2
	ECS_LOG_ERROR ECS_LOG_LEVEL = 3
	ECS_LOG_FATAL ECS_LOG_LEVEL = 4
	ECS_LOG_INFO  ECS_LOG_LEVEL = 5
)

// @title    LogLevel2Str
// @description   日志等级转字符串
// @auth      panhuili             20220826
// @param     level        ECS_LOG_LEVEL         "日志等级"
// @return     s        string         "日志等级打印字符串"
func LogLevel2Str(level ECS_LOG_LEVEL) string {
	if level == ECS_LOG_DEBUG {
		return "DEBUG"
	}
	if level == ECS_LOG_WARN {
		return "WARN "
	}
	if level == ECS_LOG_ERROR {
		return "ERROR"
	}
	if level == ECS_LOG_FATAL {
		return "FATAL"
	}
	if level == ECS_LOG_INFO {
		return "INFO "
	}
	return "NONE "
}

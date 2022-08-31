package ecs

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/Jeffail/tunny"
)

// Entiy Key类型
type KeyType int64

// World, ECS基础驱动入口
type World interface {
	BaseWorld  // World 基础接口
	LogicWorld // World 业务实现接口
}

// World基础接口
type BaseWorld interface {
	ECS_W_Construct(logicworld LogicWorld, cachecopset map[int]bool) // 初始化基础数据
	ECS_W_RegisterSys(sys Sys)                                       // 注册Sys
	ECS_W_GetEntity(key KeyType) Entity                              // 获取Entiy对象
	ECS_W_CreateEntity(key KeyType) Entity                           // 创建Entiy
	ECS_W_Push(key KeyType, msg interface{}) error                   // 插入消息到队列
	ECS_W_Handle()                                                   // 处理消息
	ECS_W_Recovery()                                                 // 自动回收资源
	ECS_W_OnMsg(e Entity, payload interface{}) error                 // 处理单个消息
}

// World 业务接口
type LogicWorld interface {
	Construct()                      // 初始化
	CreateEntity(key KeyType) Entity // 创建Entiy对象

	Prepare() error       // 启动前自定义行为
	PrintECSLog(s string) // 打印ECS框架日志接口
}

// World内部队列消息
type QueueMsg struct {
	Key KeyType     // Entiy Key
	Msg interface{} // 业务定义消息
}

// World 工作者消息
type WorkerMsg struct {
	Entity Entity      // Entiy对象
	Msg    interface{} // 业务定义消息
}

// World基础接口实现
type BaseWorldV1 struct {
	EntityPool map[KeyType]Entity // Entiy池，只在Master Go程进行添加移除操作

	SysList []Sys // Sys列表

	MsgQueue    chan QueueMsg // 消息队列
	WorkQueue   chan bool     // 全局工作者队列
	CacheCopSet map[int]bool  // 缓存型Cop列表
	LogicWorld  LogicWorld    // World业务接口实现

	WorkerPool *tunny.Pool // 工作协程池
}

var _ BaseWorld = &BaseWorldV1{}

// @title    ECS_W_Construct
// @description   初始化基础信息
// @auth      panhuili             20220826
// @param     logicworld        LogicWorld         "业务World接口实现"
// @param     cachecopset        map[int]bool         "缓存型Cop列表"
func (this *BaseWorldV1) ECS_W_Construct(logicworld LogicWorld, cachecopset map[int]bool) {
	this.EntityPool = make(map[KeyType]Entity)
	this.MsgQueue = make(chan QueueMsg, GetEcsCfg().World_MaxMsgQueueLen)
	this.LogicWorld = logicworld
	this.WorkQueue = make(chan bool, GetEcsCfg().World_MaxWorkQueue)
	this.CacheCopSet = cachecopset
}

// @title    ECS_W_RegisterSys
// @description   注册Sys
// @auth      panhuili             20220826
// @param     sys        Sys         "业务Sys接口实现"
func (this *BaseWorldV1) ECS_W_RegisterSys(sys Sys) {
	this.SysList = append(this.SysList, sys)
}

// @title    ECS_W_GetEntity
// @description   获取Entity对象
// @auth      panhuili             20220826
// @param     key        KeyType         "Entity Key"
// @return     e        Entity         "Entiy对象"
func (this *BaseWorldV1) ECS_W_GetEntity(key KeyType) Entity {
	e, ok := this.EntityPool[key]
	if !ok {
		return nil
	}
	return e
}

// @title    ECS_W_CreateEntity
// @description   创建Entity对象 并插入到Entity池
// @auth      panhuili             20220826
// @param     key        KeyType         "Entity Key"
// @return     e        Entity         "Entiy对象"
func (this *BaseWorldV1) ECS_W_CreateEntity(key KeyType) Entity {
	e := this.LogicWorld.CreateEntity(key)
	if e == nil {
		return nil
	}
	this.EntityPool[key] = e
	return e
}

// @title    ECS_W_Push
// @description   向消息队列插入请求
// @auth      panhuili             20220826
// @param     key        KeyType         "Entity Key"
// @param     msg        interface{}         "业务消息"
// @return     err        error         ""
func (this *BaseWorldV1) ECS_W_Push(key KeyType, msg interface{}) error {
	if len(this.MsgQueue) >= GetEcsCfg().World_MaxMsgQueueLen {
		this.ECS_W_PrintLog(ECS_LOG_ERROR, "Entity[%v] MsgQueue is full", key)
		return ErrorNewf("MsgQueue is full")
	}
	sh.Add(Stat_DispatchMsgNum, 1)
	this.MsgQueue <- QueueMsg{key, msg}
	this.ECS_W_PrintLog(ECS_LOG_DEBUG, "Entity[%v] Push Msg ", key)
	return nil
}

// @title    ECS_W_Start
// @description   启动World Master Go程
// @auth      panhuili             20220826
// @param     done        chan bool         "关闭信号"
func (this *BaseWorldV1) ECS_W_Start(done chan bool) {
	// 拉起辅助协程
	sh.Run()    // 统计
	clock.Run() // 时钟

	// 创建工作协程池
	this.WorkerPool = tunny.NewFunc(GetEcsCfg().World_WorkerPoolSize, func(payload interface{}) interface{} {
		workermsg, ok := payload.(*WorkerMsg)
		if !ok {
			// 获取了非 WorkerMsg 参数, 不应出现，直接Panic
			this.ECS_W_PrintLog(ECS_LOG_FATAL, "WorkerPool Recv msg err")
			panic("WorkerPool Recv msg err")
			return nil
		}

		defer func() {
			// 工作队列 减一
			<-workermsg.Entity.ECS_E_GetWorkQueue()
			<-this.WorkQueue
		}()

		// 调起消息处理函数
		err := this.ECS_W_OnMsg(workermsg.Entity, workermsg.Msg)
		if err != nil {
			this.ECS_W_PrintLog(ECS_LOG_ERROR, err.Error())
		}
		return nil
	})

	// 开启Master Go程
	go func() {
		for {
			// 处理消息
			this.ECS_W_Handle()
			// 数据回收
			this.ECS_W_Recovery()

			select {
			case <-done:
				return
			default:
				continue
			}
		}
	}()
}

// @title    ECS_W_Handle
// @description   处理消息队列
// @auth      panhuili             20220826
func (this *BaseWorldV1) ECS_W_Handle() {
	sh.SetMax(Stat_MaxMsgQueueSize, len(this.MsgQueue))
	sh.Add(Stat_LoopNum, 1)
	for i := 0; i < GetEcsCfg().World_MaxMsgNumPerLoop; i++ {
		select {
		case msg := <-this.MsgQueue:
			sh.Add(Stat_ProcessMsgNum, 1)
			this.ECS_W_PrintLog(ECS_LOG_DEBUG, "Entity[%v] Get Msg", msg.Key)
			// 检查Entity是否已创建
			e := this.ECS_W_GetEntity(msg.Key)
			if e == nil {
				// 创建Entity
				e = this.ECS_W_CreateEntity(msg.Key)
				if e == nil {
					this.ECS_W_PrintLog(ECS_LOG_ERROR, "Entity[%v] Create Entity Fail", msg.Key)
					continue
				}
				sh.Add(Stat_CreateEntityNum, 1)
			}
			// 检查Entity专属工作队列是否已满
			if len(e.ECS_E_GetWorkQueue()) >= GetEcsCfg().Entity_MaxWorkerNum {
				this.ECS_W_PrintLog(ECS_LOG_ERROR, "Entity[%v] Rightnow Single WorkQueue[%v] > MaxEntityWorkQueue[%v]", e.ECS_E_GetKey(), len(e.ECS_E_GetWorkQueue()), GetEcsCfg().Entity_MaxWorkerNum)
				continue
			}
			// 检查全局工作队列是否已满
			if len(this.WorkQueue) >= GetEcsCfg().World_MaxWorkQueue {
				this.ECS_W_PrintLog(ECS_LOG_ERROR, "Rightnow World WorkQueue[%v] > MaxWorkQueue[%v]", len(this.WorkQueue), GetEcsCfg().World_MaxWorkQueue)
				continue
			}

			sh.SetMax(Stat_MaxEntityWorkQueue, len(e.ECS_E_GetWorkQueue()))
			sh.SetMax(Stat_MaxWorkQueue, len(this.WorkQueue))
			// 工作队列 加一
			this.WorkQueue <- true
			e.ECS_E_GetWorkQueue() <- true
			this.ECS_W_PrintLog(ECS_LOG_DEBUG, "Entity[%v] ProcessMsg Msg ", msg.Key)

			// 唤起工作者协程
			go this.WorkerPool.ProcessTimed(&WorkerMsg{e, msg.Msg}, time.Second*time.Duration(GetEcsCfg().World_WorkTimeOut))
		default:
			break
		}
	}
}

// @title    ECS_W_Recovery
// @description   缓存回收
// @auth      panhuili             20220826
func (this *BaseWorldV1) ECS_W_Recovery() {
	timenow := ECS_TimeNow()
	for key, e := range this.EntityPool {
		// 有工作的Entiy跳过 防止并发
		if len(e.ECS_E_GetWorkQueue()) != 0 {
			continue
		}
		lastactive, ok := e.ECS_E_GetAttr("lastactive").(int64)
		// 超过Entity_ClearTimeOut没有活跃 或 被标记删除 则从内存移除
		if !ok || timenow-lastactive > int64(GetEcsCfg().Entity_ClearTimeOut) || e.ECS_E_IsDelete() {
			// 依次移除Cop
			for copid, _ := range this.CacheCopSet {
				if cop := e.ECS_E_GetCop(copid); cop == nil {
					continue
				} else {
					cop.ECS_C_Clear()
				}
			}
			// 移除Entity
			e.ClearEntity()
			delete(this.EntityPool, key)
			sh.Add(Stat_ClearEntityNum, 1)
			continue
		}
		// 检查Cop是否需要移除
		if len(this.CacheCopSet) == 0 {
			continue
		}
		for copid, _ := range this.CacheCopSet {
			cop := e.ECS_E_GetCop(copid)
			if cop == nil {
				continue
			}
			// 置脏 且 已超过Cop_DirtyTimeOut，则强制写回内存
			if dirtytime := cop.ECS_C_GetDirtyTime(); dirtytime != 0 && dirtytime+int64(GetEcsCfg().Cop_DirtyTimeOut) < timenow {
				if err := cop.ECS_C_Flush(); err == nil {
					cop.ECS_C_SetDirtyTime(0)
				} else {
					// 写回失败，更新置脏时间等待下次写回
					cop.ECS_C_SetDirtyTime(timenow)
					this.ECS_W_PrintLog(ECS_LOG_ERROR, "Entity[%v] Cop[%v] Flush Err[%v]", e.ECS_E_GetKey(), copid, err.Error())
				}
			}
			//（没有激活  或 上次激活时间超过Cop_ClearTimeOut） 且不存在置脏，则移除内存
			if activetime := cop.ECS_C_GetActiveTime(); ((activetime != 0 && activetime+int64(GetEcsCfg().Cop_ClearTimeOut) < timenow) || activetime == 0) && cop.ECS_C_GetDirtyTime() == 0 {
				cop.ECS_C_Clear()
			}
		}
	}
}

// @title    ECS_W_OnMsg
// @description   处理单个消息
// @auth      panhuili             20220826
// @param     e        Entity         "Entity对象"
// @param     payload        interface{}         "业务定义消息"
// @param     err        error        ""
func (this *BaseWorldV1) ECS_W_OnMsg(e Entity, payload interface{}) error {
	var s Sys
	for _, sys := range this.SysList {
		if sys.RouteMatch(payload) {
			s = sys
			break
		}
	}
	if s == nil {
		sh.Add(Stat_MatchSysFailNum, 1)
		return ErrorNewf("Msg for Entity %v not match Sys", e.ECS_E_GetKey())
	}

	e.ECS_E_Lock()
	defer e.ECS_E_UnLock()

	if e.ECS_E_IsDelete() {
		return ErrorNewf("Entity %v is Deleted", e.ECS_E_GetKey())
	}

	timenow := ECS_TimeNow()
	e.ECS_E_SetAttr("lastactive", timenow)

	coplist := s.ECS_S_GetCopList()
	for _, copid := range coplist {
		if e.ECS_E_GetCop(copid) == nil {
			if e.ECS_E_CreateCop(copid) == nil {
				sh.Add(Stat_CreateEntityFailNum, 1)
				return ErrorNewf("Entity[%v] Cop[%v] Create Fail", e.ECS_E_GetKey(), copid)
			}
			sh.Add(Stat_CreateCopNum, 1)
		}
		cop := e.ECS_E_GetCop(copid)
		cop.ECS_C_SetActiveTime(timenow)

		if _, ok := this.CacheCopSet[copid]; !ok {
			defer cop.ECS_C_Clear()
			defer sh.Add(Stat_ClearCopNum, 1)
		}
	}
	sh.Add(Stat_RunSysMainNum, 1)
	s.MainFunc(payload, e)

	return nil
}

// @title    ECS_W_PrintLog
// @description   打印ECS框架日志
// @auth      panhuili             20220826
// @param     level        ECS_LOG_LEVEL         "日志等级"
// @param     log        string         "日志内容"
func (this *BaseWorldV1) ECS_W_PrintLog(level ECS_LOG_LEVEL, log string, v ...interface{}) {
	if level < GetEcsCfg().World_LogLevel {
		return
	}
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		panic("Could not get context info for logger!")
	}
	funcname := runtime.FuncForPC(pc).Name()
	fn := funcname[strings.LastIndex(funcname, ".")+1:]
	this.LogicWorld.PrintECSLog(fmt.Sprintf("[ECS][%s][%s:%d %s] %s", LogLevel2Str(level), file, line, fn, fmt.Sprintf(log, v...)))
}

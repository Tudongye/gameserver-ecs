package ecs

import (
	"sync"
)

// Entity ECS的Entity
type Entity interface {
	BaseEntity  // 基础接口
	LogicEntity // 业务接口
}

// Entity 基础接口
type BaseEntity interface {
	ECS_E_Construct(logicentity LogicEntity, key KeyType) // 初始化
	ECS_E_GetKey() KeyType                                // 获取Key
	ECS_E_CreateCop(copid int) Cop                        // 创建Cop并执行Cop数据加载逻辑，调用LogicEntity.CreateCop()实现
	ECS_E_GetCop(copno int) Cop                           // 获取Cop
	ECS_E_Lock()                                          // 加锁
	ECS_E_UnLock()                                        // 解锁
	ECS_E_SetAttr(key string, value interface{})          // 设置自定义参数
	ECS_E_GetAttr(key string) interface{}                 // 获取自定义参数
	ECS_E_SetDelete()                                     // 设置删除标记
	ECS_E_IsDelete() bool                                 // 是否有删除标记

	ECS_E_GetWorkQueue() chan bool // 获取Entity当前的工作数量
}

// Entity 业务接口
type LogicEntity interface {
	CreateCop(copid int, Key KeyType) Cop // 创建Cop
	ClearEntity()                         // 移除Entity前清空操作
	InitEntity() error                    // 创建Enity后初始化操作
}

// Entity 基础接口实现
type BaseEntityV1 struct {
	Key KeyType // Entity key

	CopPool     map[int]Cop            // Cop列表
	AttrPool    map[string]interface{} // 自定义参数列表
	WorkQueue   chan bool              // 当前工作数
	LogicEntity LogicEntity            // 业务接口实现

	DeleteFlag bool       // 删除标记
	Mutex      sync.Mutex // 锁
}

var _ BaseEntity = &BaseEntityV1{}

// @title    ECS_E_Construct
// @description   初始化Entity基础数据
// @auth      panhuili             20220826
// @param     logicentity        LogicEntity         "业务Entity实现"
// @param    key        KeyType         "Entity Key"
func (this *BaseEntityV1) ECS_E_Construct(logicentity LogicEntity, key KeyType) {
	this.Key = key
	this.CopPool = make(map[int]Cop)
	this.AttrPool = make(map[string]interface{})
	this.WorkQueue = make(chan bool, GetEcsCfg().Entity_MaxWorkerNum)
	this.LogicEntity = logicentity
	this.DeleteFlag = false
}

// @title    ECS_E_GetKey
// @description   获取Entity Key
// @auth      panhuili             20220826
// @return    key        KeyType         "Entity Key"
func (this *BaseEntityV1) ECS_E_GetKey() KeyType {
	return this.Key
}

// @title    ECS_E_CreateCop
// @description   创建Cop并执行Cop数据加载逻辑，调用LogicEntity.CreateCop()实现
// @auth      panhuili             20220826
// @para    copid        int         "Cop编号 业务实现定义"
// @return    c        Cop         ""
func (this *BaseEntityV1) ECS_E_CreateCop(copid int) Cop {
	c := this.LogicEntity.CreateCop(copid, this.Key)
	if c == nil {
		return nil
	}
	if !c.ECS_C_Load() {
		c.ECS_C_Clear()
		return nil
	}
	this.CopPool[copid] = c
	return c
}

// @title    ECS_E_GetCop
// @description   获取Cop引用，返回空表示Cop未创建
// @auth      panhuili             20220826
// @para    copid        int         "Cop编号 业务实现定义"
// @return    c        Cop         ""
func (this *BaseEntityV1) ECS_E_GetCop(copno int) Cop {
	if this.CopPool[copno] == nil {
		return nil
	}
	return this.CopPool[copno]
}

// @title    ECS_E_Lock
// @description   阻塞加锁,Entiy锁,保护CopList
// @auth      panhuili             20220826
func (this *BaseEntityV1) ECS_E_Lock() {
	this.Mutex.Lock()
}

// @title    ECS_E_UnLock
// @description   阻塞解锁,Entiy锁,保护CopList
// @auth      panhuili             20220826
func (this *BaseEntityV1) ECS_E_UnLock() {
	this.Mutex.Unlock()
}

// @title    ECS_E_SetAttr
// @description   设置自定义属性
// @auth      panhuili             20220826
// @para    key        string         "属性名"
// @para    value        interface{}         ""
func (this *BaseEntityV1) ECS_E_SetAttr(key string, value interface{}) {
	this.AttrPool[key] = value
}

// @title    ECS_E_GetAttr
// @description   获取自定义属性
// @auth      panhuili             20220826
// @para    key        string         "属性名"
// @return    value        interface{}         ""
func (this *BaseEntityV1) ECS_E_GetAttr(key string) interface{} {
	return this.AttrPool[key]
}

// @title    ECS_E_GetWorkQueue
// @description   获取工作队列
// @auth      panhuili             20220826
// @return    WorkQueue        chan bool         "工作队列管道"
func (this *BaseEntityV1) ECS_E_GetWorkQueue() chan bool {
	return this.WorkQueue
}

// @title    ECS_E_IsDelete
// @description   检查是否有删除标记
// @auth      panhuili             20220826
// @return    DeleteFlag         bool         "删除标记"
func (this *BaseEntityV1) ECS_E_IsDelete() bool {
	return this.DeleteFlag
}

// @title    ECS_E_SetDelete
// @description   设置删除标记
// @auth      panhuili             20220826
func (this *BaseEntityV1) ECS_E_SetDelete() {
	this.DeleteFlag = true
}

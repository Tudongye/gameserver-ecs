package ecs

import ()

// Cop ECS的component
type Cop interface {
	BaseCop  // 基础接口
	LogicCop // 业务接口
}

// Cop 基础接口
type BaseCop interface {
	ECS_C_Construct(logiccop LogicCop, key KeyType) // 初始化Cop
	ECS_C_Load() bool                               // 加载数据, 调用Logic.LoadData()实现
	ECS_C_Clear() error                             // 清空数据, 调用Logic.ClearData()实现
	ECS_C_Flush() error                             // 强制写回数据, 调用Logic.FlushData()实现
	ECS_C_GetDirtyTime() int64                      // 获取置脏时间
	ECS_C_SetDirtyTime(int64)                       // 设置置脏时间
	ECS_C_SetActiveTime(int64)                      // 设置最后活跃时间
	ECS_C_GetActiveTime() int64                     // 获取最后活跃时间
}

// Cop 业务接口
type LogicCop interface {
	LoadData() error  // 加载数据
	ClearData() error // 清除数据
	FlushData() error // 强制写回数据
}

// BaseCopV1 Cop基础接口实现
type BaseCopV1[T any] struct {
	Key KeyType // Entity Key

	DirtyTime  int64 // 置脏时间，为0表示没有置脏
	ActiveTime int64 // 最后活跃时间

	Data     *T       // 业务数据
	LogicCop LogicCop // 业务接口实现
}

var _ BaseCop = &BaseCopV1[any]{}

// @title    ECS_C_Construct
// @description   初始化Cop基础数据
// @auth      panhuili             20220826
// @param     logiccop        LogicCop         "业务Cop实现"
// @param    key        KeyType         "Entity Key"
func (this *BaseCopV1[T]) ECS_C_Construct(logiccop LogicCop, key KeyType) {
	this.DirtyTime = 0
	this.ActiveTime = 0
	this.Data = nil
	this.LogicCop = logiccop
	this.Key = key
}

// @title    ECS_C_Clear
// @description   清空数据，调用Logic.ClearData()实现
// @auth      panhuili             20220826
// @return    err        error         ""
func (this *BaseCopV1[T]) ECS_C_Clear() error {
	if err := this.LogicCop.ClearData(); err != nil {
		return err
	}
	this.Data = nil
	this.DirtyTime = 0
	this.ActiveTime = 0
	return nil
}

// @title    ECS_C_Load
// @description   加载数据，调用Logic.LoadData()实现
// @auth      panhuili             20220826
// @return    err        error         ""
func (this *BaseCopV1[T]) ECS_C_Load() bool {
	if this.LogicCop.LoadData() != nil {
		return false
	}
	return true
}

// @title    ECS_C_Flush
// @description   强制写回数据，调用Logic.FlushData()实现
// @auth      panhuili             20220826
// @return    err        error         ""
func (this *BaseCopV1[T]) ECS_C_Flush() error {
	return this.LogicCop.FlushData()
}

// @title    ECS_C_GetDirtyTime
// @description   获取Cop置脏时间
// @auth      panhuili             20220826
// @return    timestamp        int64         "置脏时的标准时间戳"
func (this *BaseCopV1[T]) ECS_C_GetDirtyTime() int64 {
	return this.DirtyTime
}

// @title    ECS_C_SetDirtyTime
// @description   设置Cop置脏时间
// @auth      panhuili             20220826
// @param     timenow        int64         "当前标准时间戳"
func (this *BaseCopV1[T]) ECS_C_SetDirtyTime(timenow int64) {
	this.DirtyTime = timenow
}

// @title    ECS_C_GetActiveTime
// @description   获取最后活跃时间
// @auth      panhuili             20220826
// @return    timestamp        int64         "最后活跃时的标准时间戳"
func (this *BaseCopV1[T]) ECS_C_GetActiveTime() int64 {
	return this.ActiveTime
}

// @title    ECS_C_SetActiveTime
// @description   设置最后活跃时间
// @auth      panhuili             20220826
// @param     timenow        int64         "当前标准时间戳"
func (this *BaseCopV1[T]) ECS_C_SetActiveTime(timenow int64) {
	this.ActiveTime = timenow
}

// @title    SetData
// @description   设置Cop数据，供业务实现使用
// @auth      panhuili             20220826
// @param     data        *T         "泛型数据指针"
func (this *BaseCopV1[T]) SetData(data *T) {
	this.Data = data
}

// @title    SetData
// @description   获取Cop数据，供业务实现使用
// @auth      panhuili             20220826
// @return     data        *T         "泛型数据指针"
func (this *BaseCopV1[T]) GetData() *T {
	return this.Data
}

// @title    GetKey
// @description   获取Entity Key 供业务实现使用
// @auth      panhuili             20220826
// @return     key        KeyType         "Entity Key"
func (this *BaseCopV1[T]) GetKey() KeyType {
	return this.Key
}

// @title    SetDirty
// @description   置脏
// @auth      panhuili             20220826
func (this *BaseCopV1[T]) SetDirty() {
	this.ECS_C_SetDirtyTime(ECS_TimeNow())
}

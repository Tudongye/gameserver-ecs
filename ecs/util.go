package ecs

import (
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/atomic"
)

// @title    ErrorNewf
// @description   创建error 一个糖
// @auth      panhuili             20220826
// @param     format        string         "错误信息"
// @return     err        error         "Cop对象"
func ErrorNewf(format string, v ...interface{}) error {
	return errors.New(fmt.Sprintf(format, v...))
}

// 时钟对象
type Clock struct {
	Clock *atomic.Int64 // 时间戳原子缓存
}

var clock Clock

// @title    Run
// @description   运行时钟驱动协程
// @auth      panhuili             20220826
func (this *Clock) Run() {
	this.Clock = atomic.NewInt64(0)
	this.Clock.Store(time.Now().Unix())
	go func() {
		for {
			<-time.After(1 * time.Second)
			this.Clock.Store(time.Now().Unix())
		}
	}()
}

// @title    ECS_TimeNow
// @description   获取时间戳
// @auth      panhuili             20220826
// @return     timenow        int64         "时间戳"
func ECS_TimeNow() int64 {
	return clock.Clock.Load()
}

// 统计
type StatHelper struct {
	ClearTime int64 // 上次记录时间

	MetaPool []*StatUnit           // 统计数据Meta
	DataPool map[int]*atomic.Int64 // 统计数据
}

var sh StatHelper

// @title    SetMax
// @description   设置最大统计值
// @auth      panhuili             20220826
// @param     st        *StatUnit         "统计数据Meta"
// @param     d        int         "数据值"
func (this *StatHelper) SetMax(st *StatUnit, d int) {
	if d > this.Get(st) {
		this.Set(st, d)
	}
}

// @title    Set
// @description   设置统计值
// @auth      panhuili             20220826
// @param     st        *StatUnit         "统计数据Meta"
// @param     d        int         "数据值"
func (this *StatHelper) Set(st *StatUnit, d int) {
	this.DataPool[st.Code].Store(int64(d))
}

// @title    Get
// @description   获取统计值
// @auth      panhuili             20220826
// @param     st        *StatUnit         "统计数据Meta"
// @return     d        int         "数据值"
func (this *StatHelper) Get(st *StatUnit) int {
	return int(this.DataPool[st.Code].Load())
}

// @title    Add
// @description   统计值递增
// @auth      panhuili             20220826
// @param     st        *StatUnit         "统计数据Meta"
// @param     d        int         "增加值"
func (this *StatHelper) Add(st *StatUnit, d int) {
	this.Set(st, this.Get(st)+d)
}

// @title    Run
// @description   运行统计驱动协程
// @auth      panhuili             20220826
func (this *StatHelper) Run() {
	go func() {
		for {
			<-time.After(time.Duration(GetEcsCfg().World_StatPrintInter) * time.Second)
			timenow := ECS_TimeNow()
			log.Printf("ClearTime:%v TimeNow:%v Diff:%v", this.ClearTime, timenow, timenow-this.ClearTime)
			for code, st := range this.MetaPool {
				if st == nil {
					continue
				}
				if data, ok := this.DataPool[code]; ok {
					fmt.Printf("Code:%v Name:%v Data:%v Desc:%v\n", code, st.Name, data, st.Desc)
				} else {
					fmt.Printf("Code:%v Name:%v Data:nil Desc:%v\n", code, st.Name, st.Desc)
				}
				this.Set(st, 0)
			}
			this.ClearTime = timenow
		}
	}()
}

// 统计数据Meta
type StatUnit struct {
	Code int    // 编号 唯一
	Name string // 统计项名
	Desc string // 描述
}

// @title    NewStat
// @description   获取统计数据Meta
// @auth      panhuili             20220826
// @param     code        int         "编号"
// @param     name        string         "统计项名"
// @param     desc        string         "描述"
// @param     st        *StatUnit         "统计数据Meta"
func NewStat(code int, name string, desc string) *StatUnit {
	if code >= Stat_Max_Code {
		log.Fatalf("Repeat Stat %v-%v-%v", code, name, desc)
	}
	if sh.ClearTime == 0 {
		sh.ClearTime = time.Now().Unix()
		sh.MetaPool = make([]*StatUnit, Stat_Max_Code)
		sh.DataPool = make(map[int]*atomic.Int64)
	}
	if sh.MetaPool[code] != nil {
		log.Fatalf("Repeat Stat %v-%v-%v", code, name, desc)
	}
	sh.MetaPool[code] = &StatUnit{code, name, desc}
	sh.DataPool[code] = atomic.NewInt64(0)
	return sh.MetaPool[code]
}

var Stat_ProcessMsgNum *StatUnit = NewStat(1, "ProcessMsgNum", "处理消息数")
var Stat_LoopNum *StatUnit = NewStat(2, "LoopNum", "帧数")
var Stat_MaxMsgQueueSize *StatUnit = NewStat(3, "MaxMsgQueueSize", "消息队列最大消息数")
var Stat_MaxEntityWorkQueue *StatUnit = NewStat(4, "MaxEntityWorkQueue", "单个Entity最大工作者数")
var Stat_MaxWorkQueue *StatUnit = NewStat(5, "MaxWorkQueue", "最大工作者数")
var Stat_CreateEntityNum *StatUnit = NewStat(6, "CreateEntityNum", "创建Entity数")
var Stat_ClearEntityNum *StatUnit = NewStat(7, "ClearEntityNum", "回收Entity数")
var Stat_CreateCopNum *StatUnit = NewStat(8, "CreateCopNum", "创建Cop数")
var Stat_ClearCopNum *StatUnit = NewStat(9, "ClearCopNum", "回收Cop数")
var Stat_LoadCopNum *StatUnit = NewStat(10, "LoadCopNum", "加载Cop数")
var Stat_CreateEntityFailNum *StatUnit = NewStat(11, "CreateEntityFail", "创建Entity失败数")
var Stat_CreateCopFailNum *StatUnit = NewStat(12, "CreateCopFailNum", "创建Cop数失败数")
var Stat_RunSysMainNum *StatUnit = NewStat(13, "RunSysMainNum", "执行Sys主函数次数")
var Stat_MatchSysFailNum *StatUnit = NewStat(14, "MatchSysFailNum", "匹配Sys失败次数")
var Stat_DispatchMsgNum *StatUnit = NewStat(15, "DispatchMsgNum", "接受消息数")

const (
	Stat_Max_Code = 16
)

package ecs

// Sys ECS的System
type Sys interface {
	BaseSys  // 基础接口
	LogicSys // 业务接口
}

// Sys 基础接口
type BaseSys interface {
	ECS_S_Construct(logicsys LogicSys)         // 初始化基础信息
	ECS_S_RegisterCop(copid int)               // 注册Cop依赖
	ECS_S_GetCopList() []int                   // 获取依赖Cop列表
	ECS_S_GetCop(entity Entity, copid int) Cop // 从Entity提取Cop
}

// Sys 业务接口
type LogicSys interface {
	Construct()                                        // 初始化
	RouteMatch(payload interface{}) bool               // 路由匹配
	MainFunc(payload interface{}, entity Entity) error // 业务逻辑入口
}

// Sys 基础接口实现
type BaseSysV1 struct {
	CopList  []int        // 依赖Cop列表
	CopSet   map[int]bool // 依赖Cop Set 辅助检查
	LogicSys LogicSys     // 业务接口实现
}

// @title    ECS_S_Construct
// @description   初始化Sys基础数据
// @auth      panhuili             20220826
// @param     logicsys        LogicSys         "业务Sys实现"
func (this *BaseSysV1) ECS_S_Construct(logicsys LogicSys) {
	this.LogicSys = logicsys
	this.CopSet = make(map[int]bool)
}

// @title    ECS_S_RegisterCop
// @description   注册Sys依赖Cop
// @auth      panhuili             20220826
// @param     copid        int         "Cop编号，业务实现定义"
func (this *BaseSysV1) ECS_S_RegisterCop(copid int) {
	this.CopList = append(this.CopList, copid)
	this.CopSet[copid] = true
}

// @title    ECS_S_GetCopList
// @description   获取Sys依赖Cop列表
// @auth      panhuili             20220826
// @return     CopList        []int         "依赖Cop编号列表"
func (this *BaseSysV1) ECS_S_GetCopList() []int {
	return this.CopList
}

// @title    ECS_S_GetCopList
// @description   从Entity提取Cop
// @auth      panhuili             20220826
// @param     entity        Entity         "Entiy对象"
// @param     copid        int         "Cop编号"
// @return     cop        Cop         "Cop对象"
func (this *BaseSysV1) ECS_S_GetCop(entity Entity, copid int) Cop {
	if this.CopSet[copid] != true {
		return nil
	}
	return entity.ECS_E_GetCop(copid)
}

// @title    SysCopHelper[Cop]
// @description   Cop类型转换，一个代码糖
// @auth      panhuili             20220826
// @param     sys        Sys         "Sys"
// @param     e        Entity         "Entiy"
// @param     copid        int         "Cop编号"
// @return     cop        T         "Cop对象"
func SysCopHelper[T Cop](sys Sys, e Entity, copid int) T {
	cop := sys.ECS_S_GetCop(e, copid)
	if cop == nil {
		var t T
		return t
	} else {
		t, _ := cop.(T)
		return t
	}
}

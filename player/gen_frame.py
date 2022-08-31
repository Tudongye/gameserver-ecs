# -*- coding:utf-8 -*-
import xml.etree.ElementTree as ET
import os
def GenCopTemplate():
    return {
        "name":"",
        "type":"",
        "cache":""
    }
    

def GenSysTemplate():
    return {
        "name":"",
        "CopList":[]
    }

def GenEntityTemplate():
    return {
        "name":"",
        "CopList":[],
        "CopInsList":[],
        "SysList":[],
    }

CopMap = {}
SysMap = {}
Entity = GenEntityTemplate()
World = {
    "name":"",
    "package":""
}
World["package"]=""
World["name"]=""

def LoadFrame():
    print("Load frame.xml...")
    file_path = r'frame.xml'
    tree = ET.parse(file_path)
    root = tree.getroot()
    World["package"] = root.get("package")
    World["name"] = root.get("name")
    if World["package"] == "" or World["name"] == "":
        print("Package[%s] or WorldName[%s] is None"%(World["package"], World["name"]))
        exit(0)
    print("Package[%s] WorldName[%s]"%(World["package"], World["name"]))
    for node in list(root):
        if node.tag == "cop":
            if node.get("name") == "" or node.get('type') == "":
                print("Cop Err Name[%s] or Type[%s] is None"%(node.get("name"), node.get('type')))
                exit(0)
            if node.get("name") in CopMap.keys():
                print("Cop Err Name[%s] repeat"%(node.get("name")))
                exit(0)
            cop = GenCopTemplate()
            cop["name"] = node.get("name")
            cop["type"] = node.get('type')
            cop["cache"] = node.get("cache")
            CopMap[cop["name"]] = cop
            print("Add Cop Name:%s Type:%s"%(cop["name"], cop["type"]))
        elif node.tag == "sys":
            if node.get("name") == "":
                print("Sys Err Name[%s]  is None"%(cop("name")))
                exit(0)
            if node.get("name") in SysMap.keys():
                print("Sys Err Name[%s] repeat"%(node.get("name")))
                exit(0)
            sys = GenSysTemplate()
            sys["name"] = node.get("name")
            for cop in list(node.findall("cop")):
                if cop.get("name") not in CopMap.keys():
                    print("Sys[%s] Cop[%s] not define"%(sys["name"], cop.get("name")))
                    exit(0)
                if cop.get("name") in sys["CopList"]:
                    print("Sys[%s] Cop[%s] repeat"%(sys["name"], cop.get("name")))
                    exit(0)
                sys["CopList"].append(cop.get("name"))
                print("Sys Name:%s Append Cop:%s"%(sys["name"], cop.get("name")))
            SysMap[sys["name"]] = sys
            print("Add Sys Name:%s"%(sys["name"]))
        elif node.tag == "entity":
            if node.get("name") == "":
                print("Entity not define Name")
                exit(0)
            if Entity["name"] != "":
                print("Entity Repeat")
                exit(0)
            Entity["name"] = node.get("name")
            for cop in list(node.findall("cop")):
                if cop.get("name") not in CopMap.keys():
                    print("Entity[%s] Cop[%s] not define"%(Entity["name"], cop.get("name")))
                    exit(0)
                if cop.get("name") in Entity["CopList"]:
                    print("Entity[%s] Cop[%s] repeat"%(Entity["name"], cop.get("name")))
                    exit(0)
                if int(cop.get("ins")) in Entity["CopInsList"]:
                    print("Entity[%s] Cop[%s] Ins[%d] repeat"%(Entity["name"], cop.get("name"), int(cop.get("ins"))))
                    exit(0)
                
                Entity["CopList"].append(cop.get("name"))
                Entity["CopInsList"].append(int(cop.get("ins")))
                print("Entity Name:%s Append Cop:%s Ins %d"%(Entity["name"], cop.get("name"), int(cop.get("ins"))))
            for sys in list(node.findall("sys")):
                if sys.get("name") not in SysMap.keys():
                    print("Entity[%s] Sys[%s] not define"%(Entity["name"], sys.get("name")))
                    exit(0)
                if sys.get("name") in Entity["SysList"]:
                    print("Entity[%s] sys[%s] repeat"%(Entity["name"], sys.get("name")))
                    exit(0)
                Entity["SysList"].append(sys.get("name"))
                print("Entity Name:%s Append Sys:%s"%(Entity["name"], sys.get("name")))
            print("Add Entity Name:%s"%(Entity["name"]))
        else:
            print("Unknow Tag[%s]"%(node.tag))
            exit(0)

def GenCopCode(copname):
    return "%s_Cop_%s"%(Entity["name"], copname)


def GenFrame_Cop(cop):
    Template = []
    Template.append("package %s"%(World["package"]))
    Template.append("")
    Template.append("import (")
    Template.append("	\"ecs\"")
    Template.append("	\"log\"")
    Template.append(")")
    Template.append("")
    Template.append("// Cop业务接口实现")
    Template.append("type %s struct {"%(cop["name"]))
    Template.append("	ecs.BaseCopV1[%s]"%(cop["type"]))
    Template.append("}")
    Template.append("")
    Template.append("var _ ecs.Cop = &%s{}"%(cop["name"]))
    Template.append("")
    Template.append("// 创建Cop时触发")
    Template.append("func (this *%s) LoadData() error {"%(cop["name"]))
    Template.append("	// key := this.GetKey()")
    Template.append("	// data := this.GetData()")
    Template.append("	// this.SetData(nil)")
    Template.append("	log.Println(\"%s LoadData not define\")"%(cop["name"]))
    Template.append("	return nil")
    Template.append("}")
    Template.append("")
    Template.append("// 移除Cop时触发")
    Template.append("func (this *%s) ClearData() error {"%(cop["name"]))
    Template.append("	log.Println(\"%s ClearData not define\")"%(cop["name"]))
    Template.append("	return nil")
    Template.append("}")
    Template.append("")
    Template.append("// 置脏Cop，写回时触发")
    Template.append("func (this *%s) FlushData() error {"%(cop["name"]))
    Template.append("	log.Println(\"%s FlushData not define\")"%(cop["name"]))
    Template.append("	return nil")
    Template.append("}")
    Template.append("")
    return Template

def GenFrame_SysDef():
    Template = []
    Template.append("package %s"%(World["package"]))
    Template.append("")
    Template.append("import (")
    Template.append("	\"ecs\"")
    Template.append(")")
    Template.append("")
    for k in SysMap.keys():
        sys = SysMap[k]
        Template.append("")
        Template.append("type %s struct {"%(sys["name"]))
        Template.append("	ecs.BaseSysV1")
        Template.append("}")
        Template.append("")
        Template.append("var _ ecs.Sys = &%s{}"%(sys["name"]))
        Template.append("")
        Template.append("func (this *%s) Construct() {"%(sys["name"]))
        Template.append("	this.ECS_S_Construct(this)")
        for cop in sys["CopList"]:
            Template.append("	this.ECS_S_RegisterCop(int(%s))"%(GenCopCode(cop)))
        Template.append("}")
        Template.append("")
        for cop in sys["CopList"]:
            Template.append("func (this *%s) Get%s(e ecs.Entity) *%s {"%(sys["name"], cop, cop))
            Template.append("	return ecs.SysCopHelper[*%s](this, e, int(%s))"%(cop, GenCopCode(cop)))
            Template.append("}")
            Template.append("")
        Template.append("")
    Template.append("")
    return Template

def GenFrame_Sys(sys):
    Template = []
    Template.append("package %s"%(World["package"]))
    Template.append("")
    Template.append("import (")
    Template.append("	\"ecs\"")
    Template.append("	\"log\"")
    Template.append(")")
    Template.append("")
    Template.append("// 路由匹配时触发")
    Template.append("func (this *%s) RouteMatch(payload interface{}) bool {"%(sys["name"]))
    Template.append("	log.Println(\"%s RouteMatch not define\")"%(sys["name"]))
    Template.append("	return false")
    Template.append("}")
    Template.append("")
    Template.append("// 运行Sys业务逻辑时触发")
    Template.append("func (this *%s) MainFunc(payload interface{}, entity ecs.Entity) error {"%(sys["name"]))
    Template.append("	log.Println(\"%s MainFunc not define\")"%(sys["name"]))
    Template.append("	return nil")
    Template.append("}")
    Template.append("")
    return Template

def GenFrame_EntityDef():
    Template = []
    Template.append("package %s"%(World["package"]))
    Template.append("")
    Template.append("import (")
    Template.append("	\"ecs\"")
    Template.append(")")
    Template.append("")
    Template.append("type %s struct {"%(Entity["name"]))
    Template.append("	ecs.BaseEntityV1")
    Template.append("}")
    Template.append("")
    Template.append("var _ ecs.Entity = &%s{}"%(Entity["name"]))
    Template.append("")
    Template.append("type %s_Cop int"%(Entity["name"]))
    Template.append("")
    Template.append("// Cop编号")
    Template.append("const (")
    for i in range(len(Entity["CopList"])):
        Template.append("	%s %s_Cop = %d"%(GenCopCode(Entity["CopList"][i]), Entity["name"], Entity["CopInsList"][i]))
    Template.append(")")
    Template.append("")
    Template.append("var CacheCopList map[int]bool")
    Template.append("")
    Template.append("// 缓存型Cop")
    Template.append("func GetCacheCopList() map[int]bool {")
    Template.append("	if CacheCopList == nil {")
    Template.append("		CacheCopList = make(map[int]bool)")
    for i in range(len(Entity["CopList"])):
        if CopMap[Entity["CopList"][i]]["cache"] == "1":
            Template.append("		CacheCopList[int(%s)] = true"%(GenCopCode(Entity["CopList"][i])))
    Template.append("	}")
    Template.append("	return CacheCopList")
    Template.append("}")
    Template.append("")
    Template.append("func (this *%s) CreateCop(copid int, key ecs.KeyType) ecs.Cop {"%(Entity["name"]))
    Template.append("	var cop ecs.Cop")
    Template.append("	if copid < 0 {")
    Template.append("		return nil")
    for i in range(len(Entity["CopList"])):
        Template.append("	} else if copid == int(%s) {"%(GenCopCode(Entity["CopList"][i])))
        Template.append("		cop = &%s{}"%(Entity["CopList"][i]))
    Template.append("	} else {")
    Template.append("		return nil")
    Template.append("	}")
    Template.append("	cop.ECS_C_Construct(cop, key)")
    Template.append("	return cop")
    Template.append("}")
    Template.append("")
    return Template

def GenFrame_Entity():
    Template = []
    Template.append("package %s"%(World["package"]))
    Template.append("")
    Template.append("import (")
    Template.append("	\"log\"")
    Template.append(")")
    Template.append("")
    Template.append("// Entity从内存移除时触发")
    Template.append("func (this *%s) ClearEntity() {"%(Entity["name"]))
    Template.append("	log.Println(\"%s ClearEntity not define\")"%(Entity["name"]))
    Template.append("}")
    Template.append("")
    Template.append("// 创建新的Entity缓存时触发")
    Template.append("func (this *%s) InitEntity() error {"%(Entity["name"]))
    Template.append("	log.Println(\"%s InitEntity not define\")"%(Entity["name"]))
    Template.append("	return nil")
    Template.append("}")
    Template.append("")
    return Template

def GenFrame_WorldDef():
    Template = []
    Template.append("package %s"%(World["package"]))
    Template.append("")
    Template.append("import (")
    Template.append("	\"ecs\"")
    Template.append(")")
    Template.append("")
    Template.append("type %s struct {"%(World["name"]))
    Template.append("	ecs.BaseWorldV1")
    Template.append("}")
    Template.append("")
    Template.append("var _ ecs.World = &%s{}"%(World["name"]))
    Template.append("")
    Template.append("func (this *%s) CreateEntity(key ecs.KeyType) ecs.Entity {"%(World["name"]))
    Template.append("	entity := &%s{}"%(Entity["name"]))
    Template.append("	entity.ECS_E_Construct(entity, key)")
    Template.append("	entity.InitEntity()")
    Template.append("	return entity")
    Template.append("}")
    Template.append("")
    Template.append("func (this *%s) Construct() {"%(World["name"]))
    Template.append("	this.ECS_W_Construct(this, GetCacheCopList())")
    for k in SysMap.keys():
        sys = SysMap[k]
        Template.append("	%s := &%s{}"%(sys["name"],sys["name"]))
        Template.append("	%s.Construct()"%(sys["name"]))
        Template.append("	this.ECS_W_RegisterSys(%s)"%(sys["name"]))
    Template.append("}")
    Template.append("")
    Template.append("func (this *%s) Start(done chan bool) {"%(World["name"]))
    Template.append("	this.ECS_W_Start(done)")
    Template.append("}")
    Template.append("")
    Template.append("func (this *%s) PushMsg(key ecs.KeyType, msg interface{}) error {"%(World["name"]))
    Template.append("	return this.ECS_W_Push(key, msg)")
    Template.append("}")
    Template.append("")
    Template.append("func Create%s() *%s {"%(World["name"], World["name"]))
    Template.append("	world := &%s{}"%(World["name"]))
    Template.append("	world.Construct()")
    Template.append("	world.Prepare()")
    Template.append("	return world")
    Template.append("}")
    Template.append("")
    return Template


def GenFrame_World():
    Template = []
    Template.append("package %s"%(World["package"]))
    Template.append("")
    Template.append("import (")
    Template.append("	\"log\"")
    Template.append(")")
    Template.append("")
    Template.append("// World.Start时触发")
    Template.append("func (this *%s) Prepare() error {"%(World["name"]))
    Template.append("	log.Println(\"%s Prepare not define\")"%(World["name"]))
    Template.append("	return nil")
    Template.append("}")
    Template.append("")
    Template.append("// ECS框架日志接口")
    Template.append("func (this *%s) PrintECSLog(s string) {"%(World["name"]))
    Template.append("	log.Println(s)")
    Template.append("}")
    Template.append("")
    return Template


def GenFrame():
    print("GenFrame...")
    # 创建Cop文件
    for key in CopMap.keys():
        cop = CopMap[key]
        filename = "cop_%s.go"%cop["name"]
        if os.path.exists(filename):
            continue
        c = GenFrame_Cop(cop)
        f = open(filename, "w")
        for l in c:
            f.write("%s\n"%l)
        print("Create: %s"%filename)

    # 创建 SysDef
    filename = "sys_def.go"
    c = GenFrame_SysDef()
    f = open(filename, "w")
    for l in c:
        f.write("%s\n"%l)
    print("Update: %s"%filename)

    # 创建 Sys文件
    for key in SysMap.keys():
        sys = SysMap[key]
        filename = "sys_%s.go"%sys["name"]
        if os.path.exists(filename):
            continue
        c = GenFrame_Sys(sys)
        f = open(filename, "w")
        for l in c:
            f.write("%s\n"%l)
        print("Create: %s"%filename)
    
    # 创建 EntityDef
    filename = "entity_def.go"
    c = GenFrame_EntityDef()
    f = open(filename, "w")
    for l in c:
        f.write("%s\n"%l)
    print("Update: %s"%filename)

    # 创建 Entity文件
    while True:
        filename = "entity_%s.go"%(Entity["name"])
        if os.path.exists(filename):
            break
        c = GenFrame_Entity()
        f = open(filename, "w")
        for l in c:
            f.write("%s\n"%l)
        print("Create: %s"%filename)
        break

    # 创建 WorldDef
    filename = "world_def.go"
    c = GenFrame_WorldDef()
    f = open(filename, "w")
    for l in c:
        f.write("%s\n"%l)
    print("Update: %s"%filename)

    
    # 创建 World文件
    while True:
        filename = "world_%s.go"%(World["name"])
        if os.path.exists(filename):
            break
        c = GenFrame_World()
        f = open(filename, "w")
        for l in c:
            f.write("%s\n"%l)
        print("Create: %s"%filename)
        break
LoadFrame()
GenFrame()
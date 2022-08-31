# 这是一个基于ECS模型，设计的后端开发框架。

# 框架的能力详见 



# 代码介绍

ecs // 抽象的基础库，包含Entity（实体），Component（组件），System（系统），World（框架驱动）的接口定义 和 基础功能实现
- cop.go // Component（组件）
- ecs.go // 默认框架配置
- entity.go // Entity（实体）
- sys.go // System（系统）
- util.go // 一些日志和统计信息相关的工具
- world.go // 框架驱动，包含大部分的驱动逻辑

player  // 一个样例模块，模块功能详见前述文章
- cop_BagCop.go // 脚手架生成，背包组件
- cop_PlayerCacheCop.go // 脚手架生成，玩家缓存组件
- cop_WeaponCop.go // 脚手架生成，武器组件
- data.go // 自定义业务数据结构
- entity_def.go // 脚手架生成，和组件相关的辅助接口（完全由生成工具根据框架定义文件生成）
- entity_PlayerEntity.go // 脚手架生成，实体
- frame.xml // 框架定义文件
- gen_frame.py // 脚手架代码生成工具
- sys_ActiveWeaponSys.go // 脚手架生成，激活武器系统
- sys_CreatePlayerSys.go // 脚手架生成，创建玩家系统
- sys_def.go // 脚手架生成，和系统相关的辅助接口（完全由生成工具根据框架定义文件生成）
- sys_GetWeaponSys.go // 脚手架生成，获取武器列表系统
- world_def.go  // 脚手架生成，和引擎相关的辅助接口（完全由生成工具根据框架定义文件生成）
- world_PlayerWorld.go  // 脚手架生成，引擎接口

main // 样例测试入口
- main.go // 

![1661858126-857-630df14ed13fe-408906](https://user-images.githubusercontent.com/16680818/187608391-23177bf7-c5a5-461a-9059-221d986d7f58.png)


# 使用方法(以上述player模块为例)

1、拷贝frame.xml 和 gen_frame.py 到新目录player

2、根据需求修改frame.xml

3、执行 python gen_frame.py 生成脚手架代码

4、创建data.go 定义业务数据结构

5、根据业务逻辑修改cop_*.go 和 sys_*.go

6、编写测试代码，在go.1.18下执行(因为有使用泛型)


# 测试代码输出

```
2022/08/30 04:22:52 PlayerWorld Prepare not define      // 环境初始化
2022/08/30 04:22:52 Call: {ActiveWeaponSys 1 1001 给玩家1激活武器1001}  // 发起第一次请求 激活武器系统(依赖武器组件 和 背包组件)
2022/08/30 04:22:52 PlayerEntity InitEntity 
2022/08/30 04:22:52 WeaponCop ClearData
2022/08/30 04:22:52 [ECS][ERROR][/code/frame/ecs/world.go:163 func1] Entity[1] Cop[1] Create Fail   // 武器组件初始化失败，第一次请求失败
2022/08/30 04:22:53 Call: {CreatePlayerSys 1 0 创建玩家1}   // 发起第二次请求  创建玩家系统（依赖缓存组件）
2022/08/30 04:22:53 CreatePlayerSys MainFunc    // 进入创建玩家系统业务接口，系统调用成功，创建武器组件 和 背包组件
2022/08/30 04:22:54 Call: {ActiveWeaponSys 1 1001 给玩家1激活武器1001}     // 发起第三次请求 激活武器系统(依赖武器组件 和 背包组件)
2022/08/30 04:22:54 ActiveWeaponSys MainFunc    // 进入激活武器系统业务接口
2022/08/30 04:22:54 UseGold 5 , Left: 95        // 消耗5金币
2022/08/30 04:22:54 ActiviteWeapon: 1001        // 激活武器1001
2022/08/30 04:22:55 Call: {GetWeaponSys 1 0 获取玩家1得武器列表}    // 发起第四次请求 获取武器列表系统(依赖武器组件)
2022/08/30 04:22:55 GetWeaponSys MainFunc   // 进入获取武器列表系统业务接口
2022/08/30 04:22:55 GetWeaponList: [1001]   // 打印玩家所拥有的武器
2022/08/30 04:23:05 WeaponCop ClearData     // 超过10秒（配置）没有使用，释放武器组件
2022/08/30 04:23:05 BagCop ClearData        // 超过10秒（配置）没有使用，释放背包组件
```



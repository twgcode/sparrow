# sparrow

sparrow 是一个基于 `gin` 的 restful api 风格的 web 框架

## 框架特色

原生兼容 `gin` 框架，该框架主要减少日常开发时搭框架的工作量; sparrow 内部提供了大量的开箱即用的功能，避免开发进行复杂来实现各个组件的配合。<br/>
sparrow 在提供减少配置的复杂性时,会保证避免对使用的开源组件的源码修改，从而保证框架可以使用各个开源组件的最新版本。<br/>
sparrow 注重提供减少开发者的工作量, 但是不限制开发者编码的自由, 也就是说sparrow会提供便利的配置,但是用不用你自己决定。

## 功能列表
- 快速生成项目结构(restful 风格)
- 对系统信号处理的支持
  - [ ] 支持http server优化关闭（通过命令行操作）
  - [ ] 支持http server优雅重启（通过命令行操作）
  - ...

- 命令行参数支持
  - [x] 使用 `cobra` 实现命令行参数支持
  - [ ] 支持查看当前http连接数

- 常用中间件
  - [x] 对panic和web访问日志记录写到框架里面内置中间件支持
  - [ ] 对其他中间件以第三方包支持
  

- 日志
  - [x] 使用 `zap` 实现日志支持
  - [x] 使用 `lumberjack` 对日志进行切割管理


- 配置读取
  - [x] 支持 code 和 file 2种方式,配置框架
  - [x] 配置文件方面使用 viper实现支持

- 常用组件支持(以第三方包的形式发布)
  - [x] mysql 使用 `sqlx` 和 `go-sql-driver` 作为mysql支持
  - redis 
  - etcd
  - elasticsearch
  - websocket
  - ...



**sparrow的目标就是**:尽量平衡自由和规范的界限，既可以让开发者减少开发工作量和开发心智负担，同时开发者有自由组合的权利；当然这是个美好的愿望，希望能实现. 😁😁😁



### 示例程序
main.go 代码如下:
```go
/**
@Author: wei-g
@Date:   2020/6/19 5:33 下午
@Description:
*/

package main

import (
	"fmt"
	"net/http"

	"github.com/twgcode/sparrow/framework"
	"github.com/twgcode/sparrow/util/data"
	"github.com/twgcode/sparrow/util/log"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var (
	err error
	mgr = &MGR{}
)

type MGR struct {
	Name string
}

func run(cmd *cobra.Command, args []string) (err error) {
	framework.UseDefaultMiddleware(false)
	framework.Sparrow.Engine.GET("/", RootHandle)
	framework.Engine.GET("/1", RootHandle)
	framework.Engine.GET("/p", PanicHandle)
	fmt.Println(mgr.Name)
	return
}

func main() {
	cfg := framework.CallSparrowCfg{
		Use:       "examples",
		Short:     "sparrow 示例项目",
		Long:      "sparrow 示例项目, 用来演示 sparrow的新特性 🎉 🎉 🎉",
		Version:   "v0.0.1",
		CallerRun: run,
		CmdCfg:    true,
		CallOnConfigChange: func(e fsnotify.Event) {
			// 配置文件发生变更之后会调用的回调函数
			framework.Sparrow.ConfigConf.Unmarshal(mgr)
			fmt.Println("Config file changed:", mgr.Name)
			log.Info(mgr.Name)
		},
		CallRawVal:              mgr,
		CallDecoderConfigOption: nil,
		CfgType:                 framework.FileType,
		SparrowCfg:              nil,
	}
	if err = framework.ConfigApp(&cfg); err != nil {
		panic(err)
	}
	err = framework.Execute()
	fmt.Println("====== end ======")
	fmt.Println(err)

}
func RootHandle(c *gin.Context) {
	c.JSON(http.StatusOK, data.SucJson("root /"))

}

func PanicHandle(c *gin.Context) {
	panic("PanicHandle PanicHandle PanicHandle PanicHandle PanicHandle PanicHandle")
}
```

在项目目录下执行以下命令:
```shell script
# 整理项目依赖
go mod tidy
# 编译项目
go build
```
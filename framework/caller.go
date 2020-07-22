/**
@Author: wei-g
@Date:   2020/6/28 11:49 上午
@Description: 对调用方 暴露的接口
*/

package framework

import (
	"github.com/fsnotify/fsnotify"

	"github.com/twgcode/sparrow/util/conf"
	"github.com/twgcode/sparrow/util/log"
	"github.com/twgcode/sparrow/util/log/access"
)

var (
	GinCfg = &conf.GinConf{
		Addr:     ":8080",
		Mode:     "debug",
		NoRoute:  true,
		NoMethod: true,
	}
	// 本地开发使用
	SparrowCfgLocal = &conf.SparrowConf{
		Gin: GinCfg,
		Log: log.LocalBusinessLogCfg(),
	}
	// 开发环境
	SparrowCfgDev = &conf.SparrowConf{
		Gin:    GinCfg,
		Access: access.DevAccessLogCfg(),
		Log:    log.DevBusinessLogCfg(),
	}
	// 测试环境
	SparrowCfgBeta = &conf.SparrowConf{
		Gin:    GinCfg,
		Access: access.DevAccessLogCfg(),
		Log:    log.BetaBusinessLogCfg(),
	}
	// uat/pre环境
	SparrowCfgUat = &conf.SparrowConf{
		Gin:    GinCfg,
		Access: access.ProdAccessLogCfg(),
		Log:    log.UatBusinessLogCfg(),
	}
	// 发布线上使用
	SparrowCfgProd = &conf.SparrowConf{
		Gin:    GinCfg,
		Access: access.ProdAccessLogCfg(),
		Log:    log.ProdBusinessLogCfg(),
	}
)

var (
	Sparrow = NewApp()            // 框架实例方便调用方 快捷使用
	Engine  = Sparrow.newEngine() // gin引擎实例方便调用方 快捷使用
	FileCfg = &CallSparrowCfg{    // 提供默认使用 File 的方式配置框架
		CmdCfg:  true,
		CfgType: FileType,
	}
	// 本地开发  提供默认使用 Code 的方式配置框架
	CodeCfgLocal = &CallSparrowCfg{
		CmdCfg:     false,
		CfgType:    CodeType,
		SparrowCfg: SparrowCfgLocal,
	}
	// 开发环境 提供默认使用 Code 的方式配置框架
	CodeCfgDev = &CallSparrowCfg{
		CmdCfg:     false,
		CfgType:    CodeType,
		SparrowCfg: SparrowCfgDev,
	}
	// uat/pre环境
	CodeCfgBeta = &CallSparrowCfg{
		CmdCfg:     false,
		CfgType:    CodeType,
		SparrowCfg: SparrowCfgBeta,
	}
	// 测试环境
	CodeCfgUat = &CallSparrowCfg{
		CmdCfg:     false,
		CfgType:    CodeType,
		SparrowCfg: SparrowCfgUat,
	}
	// 线上环境提供默认使用 Code 的方式配置框架
	CodeCfgProd = &CallSparrowCfg{
		CmdCfg:     false,
		CfgType:    CodeType,
		SparrowCfg: SparrowCfgProd,
	}
)

// ConfigApp 配置默认实例需要的一些参数
func ConfigApp(callCfg *CallSparrowCfg) (err error) {
	err = Sparrow.ConfigApp(callCfg)
	return
}

// Execute 启动服务
func Execute() (err error) {
	err = Sparrow.Execute()
	return
}

// OnConfigChange  项目/调用方 配置文件发生变更后的回调函数
func OnConfigChange(configOnConfigChange func(e fsnotify.Event)) {

}

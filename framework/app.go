/**
@Author: wei-g
@Date:   2020/6/18 5:50 下午
@Description:
*/

package framework

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/twgcode/sparrow/util/cmd"
	"github.com/twgcode/sparrow/util/conf"
	"sync"
)

// App sparrow结构体
type App struct {
	Engine                   *gin.Engine
	initOnce                 sync.Once                            // 防止并发的初始化
	callerRun                func(*cobra.Command, []string) error // 业务调用方在web服务启动前要执行的代码
	EtcConf                  *conf.ViperConf                      //  框架配置
	ConfigConf               *conf.ViperConf                      // 项目/业务方 使用配置
	configOnConfigChange     func(e fsnotify.Event)               // 项目/业务方配置文件发生变化后的 回调函数
	configOnConfigChangeLock sync.Mutex
}

// NewApp 构造
func NewApp() (app *App) {
	app = &App{
		Engine:   gin.New(),
		initOnce: sync.Once{},
	}
	return
}

// ConfigApp 配置框架需要的一些参数
func (a *App) ConfigApp(use, short, long string, callerRun func(*cobra.Command, []string) error, configOnConfigChange func(e fsnotify.Event)) (err error) {
	a.initOnce.Do(func() {
		defer a.configOnConfigChangeLock.Unlock()
		a.callerRun = callerRun // 调用在web启动前要执行的代码
		a.configOnConfigChangeLock.Lock()
		a.configOnConfigChange = configOnConfigChange
		if err = cmd.InitCmd(use, short, long, a.runAppE); err != nil {
			return
		}
	})
	return
}

// runPre 启动web服务前的工作
func (a *App) runPre() (err error) {
	if err = a.checkCmd(); err != nil {
		return
	}
	// 初始化配置
	if err = a.initConf(); err != nil {
		return
	}
	// 初始化日志

	// 配置 Engine

	return
}

func (a *App) runAppE(cmd *cobra.Command, args []string) (err error) {
	if err = a.runPre(); err != nil {
		return
	}
	// 业务调用方执行的代码
	if a.callerRun != nil {
		if err = a.callerRun(cmd, args); err != nil {
			return
		}
	}
	// 启动 web 服务
	err = a.Engine.Run()
	return
}

// Execute 启动服务
func (a *App) Execute() (err error) {
	if err = cmd.RootCmd.Execute(); err != nil {
		return
	}
	return
}

// initConf 初始化 配置
func (a *App) initConf() (err error) {
	// 框架本身的日志

	if a.EtcConf, err = conf.NewViperConfig(cmd.GetEtc(), cmd.GetEtcEnvPrefix(), cmd.GetEtcAutoEnv()); err != nil {
		return
	}
	// 这里还差a.EtcConf的 OnConfigChange方法没 设置
	a.EtcConf.OnConfigChange(func(e fsnotify.Event) {})
	_config := cmd.TrimSpaceConfig()
	if len(_config) > 0 {
		// 项目/业务方 开启了业务配置
		if a.ConfigConf, err = conf.NewViperConfig(_config, cmd.GetConfigEnvPrefix(), cmd.GetConfigAutoEnv()); err != nil {
			return
		}
		a.ConfigConf.OnConfigChange(a.configOnConfigChange)
	}
	return
}

// OnConfigChange 配置 项目/调用方 配置文件发生变更后的回调函数
func (a *App) OnConfigChange(configOnConfigChange func(e fsnotify.Event)) {
	defer a.configOnConfigChangeLock.Unlock()
	a.configOnConfigChangeLock.Lock()
	if a.ConfigConf != nil {
		a.configOnConfigChange = configOnConfigChange
		a.ConfigConf.OnConfigChange(configOnConfigChange)
	}
}

// checkCmd 检查cmd参数格式是否正确
func (a *App) checkCmd() (err error) {
	if len(cmd.TrimSpaceEtc()) == 0 {
		err = fmt.Errorf("sparrow framework configuration file path cannot be empty, please use -e / --etc to specify the correct path")
	}
	return
}

/*
	// 抽离出来写到一个包中
	// NoRoute 设置 404 code 的 handle
	func (a *App) NoRoute() {

	}

	// NoRoute 设置 405 code 的 handle
	func (a *App) NoMethod() {

	}
*/

/*
编写config模块剩余的
编写日志模块
编写路由日志 和 recover 模块
*/

/**
@Author: wei-g
@Date:   2020/6/18 5:50 下午
@Description:
*/

package framework

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/twgcode/sparrow/util/cmd"
	"sync"
)

// App
type App struct {
	Engine     *gin.Engine
	initOnce   sync.Once // 防止并发的初始化
	runAppOnce sync.Once
}

// NewApp
func NewApp() (app *App) {
	app = &App{
		Engine:   &gin.Engine{},
		initOnce: sync.Once{},
	}
	return
}

// ConfigApp 配置框架需要的一些参数
func (a *App) ConfigApp(use, short, long string, callerRun func(*cobra.Command, []string) error) (err error) {
	a.initOnce.Do(func() {
		// 初始化命令参数
		if err = cmd.InitCmd(use, short, long, callerRun); err != nil {
			return
		}
	})
	return
}

// RunApp 让客户调用，最终写到 ConfigApp 方法的 callerRun 参数中
func (a *App) RunApp() (err error) {
	a.runAppOnce.Do(func() {
		if err = a.checkCmd(); err != nil {
			return
		}
		// 初始化配置
		// 检查
		// 配置 Engine

		// 启动 web 服务
		fmt.Println("a.Engine.Run()")
		err = a.Engine.Run()
	})
	return
}

// Execute 启动服务
func (a *App) Execute() (err error) {
	if err = cmd.RootCmd.Execute(); err != nil {
		return
	}
	return
}

// check 检查命令和配置文件等重要的选项,在 RunApp 中调用
func (a *App) check() (err error) {

	return
}

// cmd 检查cmd参数格式是否正确
func (a *App) checkCmd() (err error) {
	config := cmd.TrimSpaceConfig()
	if len(config) == 0 {
		err = fmt.Errorf("the profile path cannot be empty, please use -c/--config to specify the correct path")
	}
	return
}

// NoRoute 设置 404 code 的 handle
func (a *App) NoRoute() {
}

// NoRoute 设置 405 code 的 handle
func (a *App) NoMethod() {

}

/*
编写config模块
编写日志模块
编写路由日志 和 recover 模块
*/

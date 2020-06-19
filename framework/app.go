/**
@Author: wei-g
@Date:   2020/6/18 5:50 下午
@Description:
*/

package framework

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/twgcode/sparrow/util/cmd"
	"sync"
)

// App
type App struct {
	Engine   *gin.Engine
	initOnce sync.Once // 防止并发的初始化
}

// NewApp
func NewApp() *App {
	return &App{
		Engine:   &gin.Engine{},
		initOnce: sync.Once{},
	}
}

// RunApp 初始化 框架
func (a *App) RunApp(use, short, long string, run func(*cobra.Command, []string)) (err error) {
	a.initOnce.Do(func() {
		// 初始化命令参数
		if err = cmd.InitCmd(use, short, long, run); err != nil {
			return
		}

	})
	return
}

// InitApp 写一些配置
func (a *App) InitApp() {

}

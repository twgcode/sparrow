/**
@Author: wei-g
@Date:   2020/6/18 5:50 下午
@Description:
*/

package framework

import (
	"github.com/gin-gonic/gin"
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

// InitApp 初始化 框架
func (a *App) InitApp() (err error) {
	a.initOnce.Do(func() {
		// 这里写初始化逻辑
	})
	return
}

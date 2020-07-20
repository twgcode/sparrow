/**
@Author: wei-g
@Date:   2020/6/28 11:49 上午
@Description: 对调用方 暴露的接口
*/

package framework

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var DefaultCfgType = FileType
var Sparrow = NewApp()

// ConfigApp 配置框架需要的一些参数
func ConfigApp(use, short, long string, callerRun func(*cobra.Command, []string) error, configOnConfigChange func(e fsnotify.Event), cfg cfgType) (err error) {
	err = Sparrow.ConfigApp(use, short, long, callerRun, configOnConfigChange, cfg)
	return
}

// Execute 启动服务
func Execute() (err error) {
	err = Sparrow.Execute()
	return
}

// OnConfigChange  项目/调用方 配置文件发生变更后的回调函数
func OnConfigChange(configOnConfigChange func(e fsnotify.Event)) {
	Sparrow.OnConfigChange(configOnConfigChange)

}

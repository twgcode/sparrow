/**
@Author: wei-g
@Date:   2020/6/18 5:50 下午
@Description:
*/

package framework

import (
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/twgcode/sparrow/handle"
	"github.com/twgcode/sparrow/util/cmd"
	"github.com/twgcode/sparrow/util/conf"
	"github.com/twgcode/sparrow/util/log"
	"github.com/twgcode/sparrow/util/log/access"
)

type cfgType int8

const (
	FileType cfgType = 0
	CodeType cfgType = 1
)

var (
	cfgTypeErr = fmt.Errorf("cfgType wrong value")
)

// CallSparrowCfg 调用框架时需要的配置
type CallSparrowCfg struct {
	Use            string
	Short          string
	Long           string
	Version        string
	AssertValidate bool // 是否 默认调用  AssertValidate()
	CallerRun      func(*cobra.Command, []string) error
	CmdCfg         bool // 控制 项目本身配置是否由cmd指定的 -c 参数指定
	// 如果 CmdCfg 为 true 下面的这 3 个值才有意义
	CallRawVal              interface{} // 调用方 配置结构体 实例, 一定要是指针类型
	CallOnConfigChange      func(e fsnotify.Event)
	CallDecoderConfigOption []viper.DecoderConfigOption

	CfgType    cfgType // 框架读取框架配置的方式: file, code
	SparrowCfg *conf.SparrowConf
}

// App sparrow结构体
type App struct {
	Engine         *gin.Engine
	callSparrowCfg *CallSparrowCfg

	EtcConf                  *conf.ViperConf //  框架配置
	ConfigConf               *conf.ViperConf // 项目/业务方 使用配置
	configOnConfigChangeLock sync.Mutex
	initOnce                 sync.Once // 防止多次调用
	newGinOnce               sync.Once // 防止多次调用
}

// NewApp 构造
func NewApp() (app *App) {
	app = &App{}
	return
}

// ConfigApp 配置框架需要的一些参数
func (a *App) ConfigApp(callCfg *CallSparrowCfg) (err error) {
	a.initOnce.Do(func() {
		a.callSparrowCfg = callCfg
		if err = cmd.InitCmd(callCfg.Use, callCfg.Short, callCfg.Long, a.runAppE, callCfg.Version); err != nil {
			return
		}
	})
	return
}

// setGin 设置 gin 有关配置
func (a *App) setGin() {
	// 设置 gin 的 模式
	gin.SetMode(a.callSparrowCfg.SparrowCfg.Gin.Mode)
	// 添加一些 默认的 handle
	if a.callSparrowCfg.SparrowCfg.Gin.NoRoute {
		a.Engine.NoRoute(handle.NoRoute)
	}
	if a.callSparrowCfg.SparrowCfg.Gin.NoMethod {
		a.Engine.NoMethod(handle.NoMethod)
	}
}

// runPre 启动web服务前的工作
func (a *App) runPre() (err error) {
	// 检查命令传参
	if err = a.checkCmd(); err != nil {
		return
	}
	// 1 初始化配置
	if err = a.initConf(); err != nil {
		return
	}
	// 2 初始化日志
	// 2-1 构建业务日志
	if _, _, err = log.NewBusinessLoggerMgr(a.callSparrowCfg.SparrowCfg.Log); err != nil {
		return
	}
	// 2-2 构建路由日志
	if a.callSparrowCfg.SparrowCfg.Access != nil {
		if _, _, err = access.NewLoggerMgr(a.callSparrowCfg.SparrowCfg.Access); err != nil {
			log.Error("无法加载 路由日志", zap.Error(err))
			return
		}
	}

	// AssertValidate()
	if a.callSparrowCfg.AssertValidate {
		if _, err = AssertValidate(); err != nil {
			log.Error("无法快速使用框架中自带的Validator", zap.Error(err))
			return
		}
		log.Info("gin Validator.Engine type assertion is *validator.Validate success")
	}

	// 配置 Engine
	a.setGin()
	return
}

func (a *App) runAppE(cmd *cobra.Command, args []string) (err error) {
	a.newEngine()
	if err = a.runPre(); err != nil {
		return
	}
	// 业务调用方执行的代码
	if a.callSparrowCfg.CallerRun != nil {
		if err = a.callSparrowCfg.CallerRun(cmd, args); err != nil {
			return
		}
	}
	// 启动 web 服务
	err = a.Engine.Run(a.callSparrowCfg.SparrowCfg.Gin.Addr)
	return
}
func (a *App) newEngine() *gin.Engine {
	a.newGinOnce.Do(func() {
		a.Engine = gin.New()
	})
	return a.Engine
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
	switch a.callSparrowCfg.CfgType {
	default:
		err = cfgTypeErr
		return
	case CodeType:
		if err = a.initCallConf(); err != nil {
			return
		}
	case FileType:
		if err = a.initFileConf(); err != nil {
			return
		}
		if err = a.initCallConf(); err != nil {
			return
		}
	}
	return

}

// initFileConf 以 FileType 获取gin有关配置
func (a *App) initFileConf() (err error) {
	if a.EtcConf, err = conf.NewViperConfig(cmd.GetEtc(), cmd.GetEtcEnvPrefix(), cmd.GetEtcAutoEnv()); err != nil {
		return
	}
	if err = a.EtcConf.Viper.Unmarshal(&a.callSparrowCfg.SparrowCfg); err != nil {
		return
	}
	return
}

// initCallConf 业务方配置是否需要从 cmd 指定的配置文件中读取配置
func (a *App) initCallConf() (err error) {
	if a.callSparrowCfg.CmdCfg {
		_config := cmd.TrimSpaceConfig()
		if len(_config) > 0 {
			// 项目/业务方 开启了业务配置
			if a.ConfigConf, err = conf.NewViperConfig(_config, cmd.GetConfigEnvPrefix(), cmd.GetConfigAutoEnv()); err != nil {
				return
			}
			if err = a.ConfigConf.Unmarshal(a.callSparrowCfg.CallRawVal, a.callSparrowCfg.CallDecoderConfigOption...); err != nil {
				return
			}
			// 开启 配置文件 监听
			if a.callSparrowCfg.CallOnConfigChange != nil {
				a.ConfigConf.WatchConfig()
				a.ConfigConf.OnConfigChange(a.callSparrowCfg.CallOnConfigChange)
			}
		}
	}
	return
}

// checkCmd 检查cmd参数格式是否正确
func (a *App) checkCmd() (err error) {
	if a.callSparrowCfg.CfgType == FileType {
		if len(cmd.TrimSpaceEtc()) == 0 {
			err = fmt.Errorf("sparrow framework configuration file path cannot be empty, please use -e / --etc to specify the correct path")
		}
	}
	return
}

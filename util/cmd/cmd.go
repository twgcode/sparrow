/**
@Author: wei-g
@Date:   2020/6/18 6:21 下午
@Description:
*/

package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	// 保证  InitCmd 方法只执行一次
	once                          sync.Once
	config                        string // 项目/业务方 使用配置
	etc                           string // 框架本身需要的配置
	etcAutoEnv, configAutoEnv     bool
	etcEnvPrefix, configEnvPrefix string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   filepath.Base(os.Args[0]),                    // 默认的项目名(当前程序名)
	Short: "sparrow 是一个基于 gin 的 restful api 风格的 web 框架", // 默认的项目介绍
	Long: `sparrow 是一个简洁优雅的 go 语言 web 框架,用来快速开发 前后端分离 模式下的后端业务逻辑.
sparrow的目标就是: 尽量平衡自由和规范的界限，既可以让开发者减少开发工作量和开发心智负担，同时开发者有自由组合的权利；当然这是个美好的愿望，希望能实现. 😁😁😁`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

func init() {
	// 设置 RootCmd 的Flags, 设置 Flags 一定要保证在 InitCmd() 执行前执行。
	// 在 cobra 中同一个命令的同一个Flag不能被重复添加,这里使用 init 函数执行特性之一进行保证(如果某个包被导入了多次，也只会执行一次包的初始化)
	func() {
		RootCmd.PersistentFlags().StringVarP(&etc, "etc", "e", "./etc.json", "sparrow config file")
		RootCmd.PersistentFlags().StringVarP(&config, "config", "c", "./config.json", "config file ")

		RootCmd.PersistentFlags().BoolVar(&etcAutoEnv, "etc_auto_env", true, "automaticEnv has etc check ENV variables for all .  keys set in config, default & flags from etc")
		RootCmd.PersistentFlags().BoolVar(&configAutoEnv, "config_auto_env", true, "automaticEnv has etc check ENV variables for all .  keys set in config, default & flags from config")

		RootCmd.PersistentFlags().StringVar(&etcEnvPrefix, "etc_env_prefix", "sparrow", "defines a prefix that ENVIRONMENT variables will use from etc")
		RootCmd.PersistentFlags().StringVar(&configEnvPrefix, "config_env_prefix", "config", "defines a prefix that ENVIRONMENT variables will use from config")

	}()
}

func TrimSpaceEtc() string {
	etc = strings.TrimSpace(etc)
	return etc
}
func GetEtc() string {
	return etc
}

func GetConfig() string {
	return config
}

func TrimSpaceConfig() string {
	etc = strings.TrimSpace(config)
	return config
}

func GetEtcAutoEnv() bool {
	return etcAutoEnv
}

func GetConfigAutoEnv() bool {
	return configAutoEnv
}

func GetEtcEnvPrefix() string {
	return etcEnvPrefix
}

func GetConfigEnvPrefix() string {
	return configEnvPrefix
}

// InitCmd 初始化命令行
func InitCmd(use, short, long string, runE func(*cobra.Command, []string) error) (err error) {
	once.Do(func() {
		// 初始化 RootCmd 配置
		_initRootCmd(use, short, long, runE)
		/*
			if err = RootCmd.Execute(); err != nil {
				return
			}
		*/
	})
	return
}

// _initRootCmd 初始化 RootCmd 配置,
func _initRootCmd(use, short, long string, runE func(*cobra.Command, []string) error) {
	if use != "" {
		RootCmd.Use = use
	}
	if short != "" {
		RootCmd.Short = short
	}
	if long != "" {
		RootCmd.Long = long
	}
	if runE != nil {
		RootCmd.RunE = runE
	}
}

// AddCommand 添加子命令
func AddCommand(cmd ...*cobra.Command) {
	RootCmd.AddCommand(cmd...)
}

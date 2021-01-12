/**
@Author: wei-g
@Date:   2020/6/18 6:21 下午
@Description:
*/

package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

var (
	// 保证  InitCmd 方法只执行一次
	once sync.Once
	// RootCmdFlag  快速获取 RootCmd 的各个flag
	RootCmdFlag rootCmdFlag
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   filepath.Base(os.Args[0]),                    // 默认的项目名(当前程序名)
	Short: "sparrow 是一个基于 gin 的 restful api 风格的 web 框架", // 默认的项目介绍
	Long: `sparrow 是一个简洁优雅的 go 语言 web 框架,用来快速开发 前后端分离 模式下的后端业务逻辑.
sparrow的目标就是: 尽量平衡自由和规范的界限，既可以让开发者减少开发工作量和开发心智负担，同时开发者有自由组合的权利；当然这是个美好的愿望，希望能实现. 😁😁😁`,
}

type rootCmdFlag struct {
	config                        string // 项目/业务方 使用配置
	etc                           string // 框架本身需要的配置
	etcAutoEnv, configAutoEnv     bool
	etcEnvPrefix, configEnvPrefix string
}

func (r *rootCmdFlag) Etc() string             { return strings.TrimSpace(r.etc) }
func (r *rootCmdFlag) EtcEnvPrefix() string    { return strings.TrimSpace(r.etcEnvPrefix) }
func (r *rootCmdFlag) EtcAutoEnv() bool        { return r.etcAutoEnv }
func (r *rootCmdFlag) Config() string          { return strings.TrimSpace(r.config) }
func (r *rootCmdFlag) ConfigAutoEnv() bool     { return r.configAutoEnv }
func (r *rootCmdFlag) ConfigEnvPrefix() string { return strings.TrimSpace(r.configEnvPrefix) }

func init() {
	// 设置 RootCmd 的Flags, 设置 Flags 一定要保证在 InitCmd() 执行前执行。
	// 在 cobra 中同一个命令的同一个Flag不能被重复添加,这里使用 init 函数执行特性之一进行保证(如果某个包被导入了多次，也只会执行一次包的初始化)
	func() {
		RootCmd.PersistentFlags().StringVarP(&RootCmdFlag.etc, "etc", "e", "./etc.json", "sparrow config file")
		RootCmd.PersistentFlags().StringVarP(&RootCmdFlag.config, "config", "c", "", "config file")

		RootCmd.PersistentFlags().BoolVar(&RootCmdFlag.etcAutoEnv, "etc_auto_env", true, "automaticEnv has etc check ENV variables for all .  keys set in config, default & flags from etc")
		RootCmd.PersistentFlags().BoolVar(&RootCmdFlag.configAutoEnv, "config_auto_env", true, "automaticEnv has etc check ENV variables for all .  keys set in config, default & flags from config")

		RootCmd.PersistentFlags().StringVar(&RootCmdFlag.etcEnvPrefix, "etc_env_prefix", "sparrow", "defines a prefix that ENVIRONMENT variables will use from etc")
		RootCmd.PersistentFlags().StringVar(&RootCmdFlag.configEnvPrefix, "config_env_prefix", "config", "defines a prefix that ENVIRONMENT variables will use from config")
	}()
	// 从cmd中 获取 sparrow框架 默认需要的一些配置
	sparrowFlags()
}

// InitCmd 初始化命令行
func InitCmd(use, short, long string, runE func(*cobra.Command, []string) error, version string) (err error) {
	once.Do(func() {
		// 初始化 RootCmd 配置
		_initRootCmd(use, short, long, runE, version)
	})
	return
}

// _initRootCmd 初始化 RootCmd 配置,
func _initRootCmd(use, short, long string, runE func(*cobra.Command, []string) error, version string) {
	if use != "" {
		RootCmd.Use = use
	}
	if short != "" {
		RootCmd.Short = short
	}
	if long != "" {
		RootCmd.Long = long
	}
	if version != "" {
		RootCmd.Version = version
	}
	if runE != nil {
		RootCmd.RunE = runE
	}
}

// AddCommand 添加子命令
func AddCommand(cmd ...*cobra.Command) {
	RootCmd.AddCommand(cmd...)
}

/**
@Author: wei-g
@Date:   2020/6/18 6:21 下午
@Description:
*/

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"sync"
)

var (
	AppName string = filepath.Base(os.Args[0])
	Short   string = "sparrow 是一个基于 gin 的 restful api 风格的 web 框架"
	Long    string = `sparrow 是一个简洁优雅的 go 语言 web 框架,用来快速开发 前后端分离 模式下的后端业务逻辑.
sparrow的目标就是: 尽量平衡自由和规范的界限，既可以让开发者减少开发工作量和开发心智负担，同时开发者有自由组合的权利；当然这是个美好的愿望，希望能实现. 😁😁😁`
	once sync.Once
)

type sparrowCmdData struct {
	config string // 业务配置
}

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   AppName,
	Short: Short,
	Long:  Long,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

func InitCmd(use, short, long string) {
	once.Do(func() {
		//
		RootCmd.Use = use
		RootCmd.Short = short
		RootCmd.Long = long
		if err := RootCmd.Execute(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	})
}

// AddCommand 添加子命令
func AddCommand(cmd ...*cobra.Command) {
	RootCmd.AddCommand(cmd...)
}

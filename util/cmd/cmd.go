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
	Short   string = "A brief description of your application"
	Long    string = `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`

	once sync.Once
)

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   AppName,
	Short: Short,
	Long:  Long,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

func InitCmd(use, short, long string) () {
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

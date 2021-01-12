/**
@Author: wei-g
@Date:   2020/6/24 11:00 上午
@Description: 框架提供配置有关
*/

package conf

import (
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// ViperConf cmd 配置 结构体
type ViperConf struct {
	configFile   string
	autoEnv      bool   // 是否可以从环境变量中取值
	envPrefix    string // 如果从环境变量中取值时,设置变量的前缀,只有当env为true,且该值非空时，该值才生效
	Viper        *viper.Viper
	newViperOnce sync.Once
}

// NewViperConfig 生成一个新的 cmd 配置
func NewViperConfig(configFile, envPrefix string, autoEnv bool) (vc *ViperConf, err error) {
	vc = &ViperConf{
		configFile:   configFile,
		autoEnv:      autoEnv,
		envPrefix:    envPrefix,
		Viper:        viper.New(),
		newViperOnce: sync.Once{},
	}
	err = vc.init()
	return
}

// init  初始化Viper 配置
func (v *ViperConf) init() (err error) {
	v.newViperOnce.Do(func() {
		if len(v.configFile) != 0 { // 设置了文件名时
			v.Viper.SetConfigFile(v.configFile) // 设置配置文件路径
			if err = v.Viper.ReadInConfig(); err != nil {
				return
			}
		}
		// 是否自动检测环境变量的keys
		if v.autoEnv {
			if v.envPrefix = strings.TrimSpace(v.envPrefix); len(v.envPrefix) != 0 {
				v.Viper.SetEnvPrefix(v.envPrefix)
			}
			// AutomaticEnv has Viper check ENV variables for all.  keys set in conf, default & flags
			v.Viper.AutomaticEnv()
		}
	})
	return
}

// Viper 原生功能,为了方便调用更便捷的调用

/* 命令行获取参数有关 */
func (v *ViperConf) SetDefault(key string, value interface{}) {
	v.Viper.SetDefault(key, value)
}

func (v *ViperConf) BindPFlag(key string, flag *pflag.Flag) error {
	return v.Viper.BindPFlag(key, flag)
}

func (v *ViperConf) BindPFlags(flags *pflag.FlagSet) error {
	return v.Viper.BindPFlags(flags)
}

/* 文件获取参数有关 */
// OnConfigChange 配置文件发生变更后的回调函数
func (v *ViperConf) OnConfigChange(onConfigChange func(e fsnotify.Event)) {
	if onConfigChange != nil {
		v.Viper.OnConfigChange(onConfigChange)
	}
}

// WatchConfig监听配置变化
func (v *ViperConf) WatchConfig() {
	// 监听配置变化
	v.Viper.WatchConfig()
}

func (v *ViperConf) Unmarshal(rawVal interface{}, opts ...viper.DecoderConfigOption) (err error) {
	return v.Viper.Unmarshal(rawVal, opts...)
}

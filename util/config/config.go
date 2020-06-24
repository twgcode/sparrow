/**
@Author: wei-g
@Date:   2020/6/24 11:00 上午
@Description: 框架提供配置有关
*/

package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"strings"
	"sync"
)

// ViperConfig cmd 配置 结构体
type ViperConfig struct {
	configFile     string
	watchConfig    bool                   // 监听配置变化
	onConfigChange func(e fsnotify.Event) // 配置文件发生变更之后会调用的回调函数
	autoEnv        bool                   // 是否可以从环境变量中取值
	envPrefix      string                 // 如果从环境变量中取值时,设置变量的前缀,只有当env为true,且该值非空时，该值才生效
	Viper          *viper.Viper
	newViperOnce   sync.Once
}

// NewViperConfig 生成一个新的 cmd 配置
func NewViperConfig(configFile, envPrefix string, watchConfig, autoEnv bool, onConfigChange func(e fsnotify.Event)) (vc *ViperConfig, err error) {
	vc = &ViperConfig{
		configFile:     configFile,
		watchConfig:    watchConfig,
		onConfigChange: onConfigChange,
		autoEnv:        autoEnv,
		envPrefix:      envPrefix,
		Viper:          nil,
		newViperOnce:   sync.Once{},
	}
	return
}

// NewViper 构建一个 配置实例
func (v *ViperConfig) NewViper(conf *viper.Viper, err error) {
	v.newViperOnce.Do(func() {
		conf = viper.New()
		conf.SetConfigFile(v.configFile) // 设置配置文件路径
		if err = conf.ReadInConfig(); err != nil {
			return
		}
		// 是否自动检测环境变量的keys
		if v.autoEnv {
			if v.envPrefix = strings.TrimSpace(v.envPrefix); len(v.envPrefix) != 0 {
				conf.SetEnvPrefix(v.envPrefix)
			}
			// AutomaticEnv has Viper check ENV variables for all.  keys set in config, default & flags
			conf.AutomaticEnv()
		}
		// 监听配置文件变化
		if v.watchConfig {
			conf.WatchConfig() // 监听配置变化
			conf.OnConfigChange(v.onConfigChange)
		}
		v.Viper = conf
		return
	})
}

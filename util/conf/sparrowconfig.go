/**
@Author: wei-g
@Date:   2020/6/24 7:16 下午
@Description: 框架自带的配置
*/

package conf

import "github.com/twgcode/sparrow/util/log"

// SparrowConf sparrow 框架本身的配置
type SparrowConf struct {
	Gin    *GinConf          `mapstructure:"gin" json:"gin"`
	Access *log.LoggerConfig `mapstructure:"access" json:"access"` // 框架的路由日志
	Log    *log.LoggerConfig `mapstructure:"log" json:"log"`       // 业务日志
}

// 启动 gin 框架需要有关配置
type GinConf struct {
	Addr     string `mapstructure:"addr" json:"addr"`
	Mode     string `mapstructure:"mode" json:"mode"`
	NoRoute  bool   `mapstructure:"no_route" json:"no_route"`
	NoMethod bool   `mapstructure:"no_method" json:"no_method"`
}

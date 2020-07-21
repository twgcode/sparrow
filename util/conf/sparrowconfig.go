/**
@Author: wei-g
@Date:   2020/6/24 7:16 下午
@Description: 框架自带的配置
*/

package conf

import "github.com/twgcode/sparrow/util/log"

type SparrowConf struct {
	Gin    *GinConf          `mapstructure:"gin" json:"gin"`
	Access *log.LoggerConfig `mapstructure:"access" json:"access"`
	Log    *log.LoggerConfig `mapstructure:"log" json:"log"`
}

// 启动 gin 框架需要有关配置
type GinConf struct {
	Addr     string `mapstructure:"addr" json:"addr"`
	Mode     string `mapstructure:"mode" json:"mode"`
	NoRoute  bool   `mapstructure:"no_route" json:"no_route"`
	NoMethod bool   `mapstructure:"no_method" json:"no_method"`
}

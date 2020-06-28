/**
@Author: wei-g
@Date:   2020/6/24 7:16 下午
@Description: 框架自带的配置
*/

package conf

// SparrowConfig 框架自带配置结构体
type SparrowConfig struct {
	Gin ginConf // gin 框架需要的配置
}

type ginConf struct {
	Addr     string
	Mode     string
	NoRoute  bool
	NoMethod bool
}

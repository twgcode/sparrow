/**
@Author: wei-g
@Date:   2020/7/19 6:35 下午
@Description: 一些默认的日志配置,达到让使用者开箱基于的目的
*/

package log

var (
	// DefaultBusinessFileCfg  默认 业务 日志文件配置
	// 默认 不会根据时间删除旧的日志文件
	DefaultBusinessFileCfg = &LoggerFileConfig{
		FileName:   DefaultLowLevelFileName,
		MaxSize:    200,
		MaxBackups: 100,
		Compress:   true,
	}
	// ProdEnableStacktrace 生产环境 是否开启 高级别错误 堆栈信息记录
	ProdEnableStacktrace = true
	// DefaultAccessFileCfg  默认 http方法 日志文件配置
	// 默认 不会根据时间删除旧的日志文件
	DefaultAccessFileCfg = &LoggerFileConfig{
		FileName:   "access.log",
		MaxSize:    200,
		MaxBackups: 100,
		Compress:   true,
	}
)

// LocalBusinessLogCfg 本地开发 环境 默认日志配置; 输出到控制台
func LocalBusinessLogCfg() (cfg *LoggerConfig) {
	cfg = &LoggerConfig{
		OutputConsole:       true,
		OutputFile:          false,
		SplitWriteFromLevel: false,
		Stacktrace:          true,
		LowLevel:            "debug",
		HighLevel:           "warn",
		EncoderText:         "console",
		EncoderConfigText:   EncoderConfigTextDevDefault,
	}
	return

}

// DevBusinessLogCfg 开发 环境 默认日志配置; 输出到文件
func DevBusinessLogCfg() (cfg *LoggerConfig) {
	cfg = &LoggerConfig{
		OutputConsole:       false,
		OutputFile:          true,
		SplitWriteFromLevel: false,
		Stacktrace:          true,
		LowLevel:            "debug",
		HighLevel:           "warn",
		LowLevelFile:        DefaultBusinessFileCfg,
		EncoderText:         "console",
		EncoderConfigText:   EncoderConfigTextDevDefault,
	}
	return
}

// BetaBusinessLogCfg 测试 环境 默认日志配置; 输出到文件
func BetaBusinessLogCfg() (cfg *LoggerConfig) {
	cfg = &LoggerConfig{
		OutputConsole:       false,
		OutputFile:          true,
		SplitWriteFromLevel: false,
		Stacktrace:          true,
		LowLevel:            "debug",
		HighLevel:           "error",
		LowLevelFile:        DefaultBusinessFileCfg,
		EncoderText:         "console",
		EncoderConfigText:   EncoderConfigTextDevDefault,
	}
	return
}

// UatBusinessLogCfg  uat/pre 环境 默认日志配置
// 输出到文件, 进行日志分级输出
func UatBusinessLogCfg() (cfg *LoggerConfig) {
	cfg = &LoggerConfig{
		OutputConsole:       false,
		OutputFile:          true,
		SplitWriteFromLevel: true,
		Stacktrace:          true,
		LowLevel:            "debug",
		HighLevel:           "error",
		LowLevelFile:        DefaultBusinessFileCfg,
		EncoderText:         "json",
		EncoderConfigText:   EncoderConfigTextProdDefault,
	}
	return
}

// ProdBusinessLogCfg  生产 环境 默认日志配置
// 输出到文件, 进行日志分级输出
func ProdBusinessLogCfg() (cfg *LoggerConfig) {
	cfg = &LoggerConfig{
		OutputConsole:       false,
		OutputFile:          true,
		SplitWriteFromLevel: true,
		Stacktrace:          ProdEnableStacktrace,
		LowLevel:            "info",
		HighLevel:           "error",
		LowLevelFile:        DefaultBusinessFileCfg,
		EncoderText:         "json",
		EncoderConfigText:   EncoderConfigTextProdDefault,
	}
	return
}

// ProdAccessLogCfg  生产 环境 http 访问日志默认配置
// 输出到文件, 进行日志分级输出
func DevAccessLogCfg() (cfg *LoggerConfig) {
	cfg = &LoggerConfig{
		OutputConsole:       false,
		OutputFile:          true,
		SplitWriteFromLevel: false,
		Stacktrace:          true,
		LowLevel:            "info",
		HighLevel:           "error",
		LowLevelFile:        DefaultAccessFileCfg,
		EncoderText:         "console",
		EncoderConfigText:   EncoderConfigTextDevDefault,
	}
	return
}

// ProdAccessLogCfg  生产 环境 http 访问日志默认配置
// 输出到文件, 进行日志分级输出
func ProdAccessLogCfg() (cfg *LoggerConfig) {
	cfg = &LoggerConfig{
		OutputConsole:       false,
		OutputFile:          true,
		SplitWriteFromLevel: true,
		Stacktrace:          true,
		LowLevel:            "info",
		HighLevel:           "error",
		LowLevelFile:        DefaultAccessFileCfg,
		EncoderText:         "json",
		EncoderConfigText:   EncoderConfigTextProdDefault,
	}
	return
}

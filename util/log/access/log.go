/**
@Author: wei-g
@Date:   2020/7/21 4:15 下午
@Description: web 请求访问日志
*/

package access

import (
	"sync"

	"github.com/twgcode/sparrow/util/log"

	"go.uber.org/zap"
)

var (
	LoggerMgr     *log.Logger
	SugaredLogger *zap.SugaredLogger
	LoggerMgrOnce sync.Once
	// DefaultFileCfg  默认 http方法 日志文件配置 ; 默认 不会根据时间删除旧的日志文件
	DefaultFileCfg = &log.LoggerFileConfig{
		FileName:   "access.log",
		MaxSize:    200,
		MaxBackups: 100,
		Compress:   true,
	}
)

// NewBusinessLoggerMgr 构建业务日志 记录器
func NewLoggerMgr(loggerConfig *log.LoggerConfig) (*log.Logger, *zap.Logger, error) {
	var err error
	LoggerMgrOnce.Do(func() {
		// 因为对外提供了日志输出的封装所以把 跳过的调用方法的层数要+1
		if loggerConfig.CallerSkip <= 0 {
			loggerConfig.CallerSkip = 1
		} else {
			loggerConfig.CallerSkip += 1
		}
		if LoggerMgr, err = log.NewLogger(loggerConfig); err != nil {
			return
		}
		SugaredLogger = LoggerMgr.Logger.Sugar()
		return
	})
	return LoggerMgr, LoggerMgr.Logger, err
}

func Debug(msg string, fields ...zap.Field) {
	LoggerMgr.Logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	LoggerMgr.Logger.Info(msg, fields...)
}
func Warn(msg string, fields ...zap.Field) {
	LoggerMgr.Logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	LoggerMgr.Logger.Error(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	LoggerMgr.Logger.DPanic(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	LoggerMgr.Logger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	LoggerMgr.Logger.Fatal(msg, fields...)
}

func SugarDebug(args ...interface{}) {
	SugaredLogger.Debug(args...)
}

func SugarInfo(args ...interface{}) {
	SugaredLogger.Info(args...)
}

func SugarWarn(args ...interface{}) {
	SugaredLogger.Warn(args...)
}

func SugarError(args ...interface{}) {
	SugaredLogger.Error(args...)
}

func SugarDPanic(args ...interface{}) {
	SugaredLogger.DPanic(args...)
}

func SugarPanic(args ...interface{}) {
	SugaredLogger.Panic(args...)
}

func SugarFatal(args ...interface{}) {
	SugaredLogger.Fatal(args...)
}

func SugarDebugf(template string, args ...interface{}) {
	SugaredLogger.Debugf(template, args...)
}

func SugarInfof(template string, args ...interface{}) {
	SugaredLogger.Infof(template, args...)
}

func SugarWarnf(template string, args ...interface{}) {
	SugaredLogger.Warnf(template, args...)
}

func SugarErrorf(template string, args ...interface{}) {
	SugaredLogger.Errorf(template, args...)
}

func SugarDPanicf(template string, args ...interface{}) {
	SugaredLogger.DPanicf(template, args...)
}

func SugarPanicf(template string, args ...interface{}) {
	SugaredLogger.Panicf(template, args...)
}

func SugarFatalf(template string, args ...interface{}) {
	SugaredLogger.Fatalf(template, args...)
}

// DevAccessLogCfg  开发 环境 http 访问日志默认配置
// 输出到文件, 进行日志分级输出
func DevAccessLogCfg() (cfg *log.LoggerConfig) {
	cfg = &log.LoggerConfig{
		OutputConsole:       false,
		OutputFile:          true,
		SplitWriteFromLevel: false,
		Stacktrace:          true,
		LowLevel:            "info",
		HighLevel:           "error",
		LowLevelFile:        DefaultFileCfg,
		EncoderText:         "console",
		EncoderConfigText:   log.EncoderConfigTextDevDefault,
	}
	return
}

// ProdAccessLogCfg  生产 环境 http 访问日志默认配置
// 输出到文件, 进行日志分级输出
func ProdAccessLogCfg() (cfg *log.LoggerConfig) {
	cfg = &log.LoggerConfig{
		OutputConsole:       false,
		OutputFile:          true,
		SplitWriteFromLevel: true,
		Stacktrace:          true,
		LowLevel:            "info",
		HighLevel:           "error",
		LowLevelFile:        DefaultFileCfg,
		EncoderText:         "json",
		EncoderConfigText:   log.EncoderConfigTextProdDefault,
	}
	return
}

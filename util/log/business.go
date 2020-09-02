/**
@Author: wei-g
@Date:   2020/7/19 4:16 下午
@Description: 业务日志
*/

package log

import (
	"sync"

	"go.uber.org/zap"
)

var (
	BusinessLoggerMgr     *Logger
	BusinessSugaredLogger *zap.SugaredLogger
	BusinessLoggerMgrOnce sync.Once
)

// NewBusinessLoggerMgr 构建业务日志 记录器
func NewBusinessLoggerMgr(loggerConfig *LoggerConfig) (*Logger, *zap.Logger, error) {
	var err error
	BusinessLoggerMgrOnce.Do(func() {
		// 因为对外提供了日志输出的封装所以把 跳过的调用方法的层数要+1
		if loggerConfig.CallerSkip <= 0 {
			loggerConfig.CallerSkip = 1
		} else {
			loggerConfig.CallerSkip += 1
		}
		if BusinessLoggerMgr, err = NewLogger(loggerConfig); err != nil {
			return
		}
		BusinessSugaredLogger = BusinessLoggerMgr.Logger.Sugar()
		return
	})
	return BusinessLoggerMgr, BusinessLoggerMgr.Logger, err
}

func Debug(msg string, fields ...zap.Field) {
	BusinessLoggerMgr.Logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	BusinessLoggerMgr.Logger.Info(msg, fields...)
}
func Warn(msg string, fields ...zap.Field) {
	BusinessLoggerMgr.Logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	BusinessLoggerMgr.Logger.Error(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	BusinessLoggerMgr.Logger.DPanic(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	BusinessLoggerMgr.Logger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	BusinessLoggerMgr.Logger.Fatal(msg, fields...)
}

func SugarDebug(args ...interface{}) {
	BusinessSugaredLogger.Debug(args...)
}

func SugarInfo(args ...interface{}) {
	BusinessSugaredLogger.Info(args...)
}

func SugarWarn(args ...interface{}) {
	BusinessSugaredLogger.Warn(args...)
}

func SugarError(args ...interface{}) {
	BusinessSugaredLogger.Error(args...)
}

func SugarDPanic(args ...interface{}) {
	BusinessSugaredLogger.DPanic(args...)
}

func SugarPanic(args ...interface{}) {
	BusinessSugaredLogger.Panic(args...)
}

func SugarFatal(args ...interface{}) {
	BusinessSugaredLogger.Fatal(args...)
}

func SugarDebugf(template string, args ...interface{}) {
	BusinessSugaredLogger.Debugf(template, args...)
}

func SugarInfof(template string, args ...interface{}) {
	BusinessSugaredLogger.Infof(template, args...)
}

func SugarWarnf(template string, args ...interface{}) {
	BusinessSugaredLogger.Warnf(template, args...)
}

func SugarErrorf(template string, args ...interface{}) {
	BusinessSugaredLogger.Errorf(template, args...)
}

func SugarDPanicf(template string, args ...interface{}) {
	BusinessSugaredLogger.DPanicf(template, args...)
}

func SugarPanicf(template string, args ...interface{}) {
	BusinessSugaredLogger.Panicf(template, args...)
}

func SugarFatalf(template string, args ...interface{}) {
	BusinessSugaredLogger.Fatalf(template, args...)
}

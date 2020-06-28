/**
@Author: wei-g
@Date:   2020/6/28 5:29 下午
@Description:
*/

package log

import (
	"os"
	"strings"

	"git.sixents.com/sixentsdevops/templategin/util/config"
	"git.sixents.com/sixentsdevops/templategin/util/data"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/natefinch/lumberjack"
)

var Logger *zap.Logger

// InitDefaultLog 初始化默认日志配置
func InitDefaultLog() (err error) {
	if Logger != nil {
		return
	}
	Logger, err = NewLog(config.SConf.Log.Path, config.SConf.Log.MaxSize, config.SConf.Log.MaxBackups,
		config.SConf.Log.MaxAge, config.SConf.Log.Level, config.SConf.Log.Compress)
	return
}

// getLogWriter 结合日志切割工具设置日志输出接口
func getLogWriter(filename string, maxSize, maxBackup, maxAge int, compress bool) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
		Compress:   compress,
	}

	if strings.ToLower(config.SConf.Gin.Mode) == data.GinMode.Debug {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	}
	return zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberJackLogger))
}

// getEncoder 配置日志格式
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

// NewLog 构造日志
func NewLog(filename string, maxSize, maxBackup, maxAge int, level string, compress bool) (log *zap.Logger, err error) {
	writeSync := getLogWriter(filename, maxSize, maxBackup, maxAge, compress)
	encoder := getEncoder()
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(level))
	if err != nil {
		return
	}
	// 编码器配置 打印到控制台和文件 日志级别
	core := zapcore.NewCore(encoder, writeSync, l)

	log = zap.New(core, zap.AddCaller())

	return
}

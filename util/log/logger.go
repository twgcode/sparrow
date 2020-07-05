/*
@Author: wei-g
@Email:  wei_g_it@163.com
@Date:   2020/6/28 下午9:30
@Description:
*/

package log

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

var (
	// 根据不同的 LoggerConfig.EncoderText 返回不同的生成 zapcore.Encoder 的函数
	NewEncoderFuncMap = map[string]func(zapcore.EncoderConfig) zapcore.Encoder{
		EncoderTextJSON:    zapcore.NewJSONEncoder,
		EncoderTextConsole: zapcore.NewConsoleEncoder,
	}
	// 根据不同的 LoggerConfig.EncoderConfigText 返回不同的生成 zapcore.EncoderConfig 的函数
	NewEncoderConfigFuncMap = map[string]func(*Logger) (zapcore.EncoderConfig, error){
		EncoderConfigTextDevDefault:  NewEncoderConfigTextDevDefault,
		EncoderConfigTextDevCustom:   NewEncoderConfigTextDevCustom,
		EncoderConfigTextProdDefault: NewEncoderConfigTextProdDefault,
		EncoderConfigTextProdCustom:  NewEncoderConfigTextProdCustom,
	}
)

type Logger struct {
	logConfig         *LoggerConfig // 配置
	Logger            *zap.Logger   // 记录器实例
	newZapLoggerMutex sync.Mutex    // 构建 *zap.Logger时使用(NewZapLogger函数调用)
}

// NewLogger 构建有一个 日志实例
func NewLogger(loggerConfig *LoggerConfig) (log *Logger, err error) {

	log = &Logger{
		logConfig: loggerConfig,
		Logger:    nil,
	}
	if err = log.checkLoggerConfig(); err != nil {
		return
	}
	err = log.newZapLogger()
	return
}

// checkLoggerConfig 检查 要生成新的 日志记录器的配置
func (l *Logger) checkLoggerConfig() (err error) {
	if l.logConfig == nil {
		err = fmt.Errorf("LoggerConfig cannot be nil")
		return
	}
	// 检查 logConfig 字段
	err = l.logConfig.CheckLoggerConfig()
	return
}

// newZapLogger 构建 zap Logger 实例
func (l *Logger) newZapLogger() (err error) {
	defer l.newZapLoggerMutex.Unlock()
	l.newZapLoggerMutex.Lock()
	var (
		encoderConfig zapcore.EncoderConfig
		encoder       zapcore.Encoder
		cores         []zapcore.Core
	)
	// 获取  zapcore.EncoderConfig
	if encoderConfig, err = l.GetEncoderConfig(); err != nil {
		return
	}
	// 获取 zapcore.Encoder
	if encoder, err = l.GetEncoder(encoderConfig); err != nil {
		return
	}
	// 构建 zapcore.Core
	if cores, err = l.newZapCore(encoder); err != nil {
		return
	}
	// 判断是否开启高级别日志堆栈信息记录
	if l.logConfig.Stacktrace {
		level := new(zapcore.Level)
		if err = level.UnmarshalText([]byte(l.logConfig.HighLevel)); err != nil {
			return
		}
		l.Logger = zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddStacktrace(level), zap.AddCallerSkip(l.logConfig.CallerSkip))
		return
	}
	l.Logger = zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(l.logConfig.CallerSkip))
	return
}

// newZapCore 生成 zapcore.Core(是一个最小的，快速的记录器接口) , NewCore creates a Core that writes logs to a WriteSyncer
func (l *Logger) newZapCore(encoder zapcore.Encoder) (cores []zapcore.Core, err error) {
	cores = make([]zapcore.Core, 0, 3) // 最多就3个
	// 输出到控制台
	if l.logConfig.OutputConsole {
		// 获取要记录的日志等级
		level := new(zapcore.Level)
		if err = level.UnmarshalText([]byte(l.logConfig.LowLevel)); err != nil {
			return
		}
		writeSyncEr := zapcore.AddSync(os.Stdout)
		cores = append(cores, zapcore.NewCore(encoder, writeSyncEr, level))
	}
	// 输出到文件
	if l.logConfig.OutputFile {
		level := new(zapcore.Level)
		if err = level.UnmarshalText([]byte(l.logConfig.LowLevel)); err != nil {
			return
		}
		writeSyncEr := l.GetFileLogWriter(l.logConfig.LowLevelFile)
		cores = append(cores, zapcore.NewCore(encoder, writeSyncEr, level))
		// 高级别日志 单独输出
		if l.logConfig.SplitWriteFromLevel {
			level := new(zapcore.Level)
			if err = level.UnmarshalText([]byte(l.logConfig.HighLevel)); err != nil {
				return
			}
			writeSyncEr := l.GetFileLogWriter(l.logConfig.HighLevelFile)
			cores = append(cores, zapcore.NewCore(encoder, writeSyncEr, level))
		}
	}
	return
}

// GetFileLogWriter 结合日志切割工具设置日志输出接口
func (l *Logger) GetFileLogWriter(fileConfig *LoggerFileConfig) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fileConfig.FileName,
		MaxSize:    fileConfig.MaxSize,
		MaxBackups: fileConfig.MaxBackups,
		MaxAge:     fileConfig.MaxAge,
		Compress:   fileConfig.Compress,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// GetEncoderConfig 获取  配置编码器 zapcore.EncoderConfig
func (l *Logger) GetEncoderConfig() (encoderConfig zapcore.EncoderConfig, err error) {
	funcNewEncoderConfig, ok := NewEncoderConfigFuncMap[l.logConfig.EncoderConfigText]
	if !ok {
		err = optionalValuesErr(l.logConfig.EncoderConfigText, "EncoderConfigText")
	}
	encoderConfig, err = funcNewEncoderConfig(l)
	return
}

// CustomEncoderConfig 配置自动定义 配置编码器 zapcore.EncoderConfig
func (l *Logger) CustomEncoderConfig(encoderConfig *zapcore.EncoderConfig) (err error) {
	// 配置各种 编码器
	var EncodeTime = new(zapcore.TimeEncoder)
	if err = EncodeTime.UnmarshalText([]byte(l.logConfig.TimeEncoderText)); err != nil {
		return
	}
	encoderConfig.EncodeTime = *EncodeTime

	var EncodeLevel = new(zapcore.LevelEncoder)
	if err = EncodeLevel.UnmarshalText([]byte(l.logConfig.LevelEncoderText)); err != nil {
		return
	}
	encoderConfig.EncodeLevel = *EncodeLevel

	var EncodeDuration = new(zapcore.DurationEncoder)
	if err = EncodeDuration.UnmarshalText([]byte(l.logConfig.DurationEncoderText)); err != nil {
		return
	}
	encoderConfig.EncodeDuration = *EncodeDuration

	var EncodeCaller = new(zapcore.CallerEncoder)
	if err = EncodeCaller.UnmarshalText([]byte(l.logConfig.CallerEncoderText)); err != nil {
		return
	}
	encoderConfig.EncodeCaller = *EncodeCaller

	// 配置各种 key
	encoderConfig.TimeKey = l.logConfig.TimeKey
	encoderConfig.LevelKey = l.logConfig.LevelKey
	encoderConfig.CallerKey = l.logConfig.CallerKey
	encoderConfig.MessageKey = l.logConfig.MessageKey
	encoderConfig.StacktraceKey = l.logConfig.StacktraceKey
	return
}

// GetEncoder 获取 日志编码器 zapcore.Encoder
func (l *Logger) GetEncoder(encoderConfig zapcore.EncoderConfig) (encoder zapcore.Encoder, err error) {
	// 通过 表驱动 方法实现
	funcNewEncoder, ok := NewEncoderFuncMap[l.logConfig.EncoderText]
	if !ok {
		err = optionalValuesErr(l.logConfig.EncoderText, "EncoderText")
	}
	encoder = funcNewEncoder(encoderConfig)
	return
}

// NewEncoderConfigTextProdDefault   生产环境默认的 EncoderConfig 配置
func NewEncoderConfigTextProdDefault(l *Logger) (encoderConfig zapcore.EncoderConfig, err error) {
	encoderConfig = zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return
}

func NewEncoderConfigTextProdCustom(l *Logger) (encoderConfig zapcore.EncoderConfig, err error) {
	encoderConfig = zap.NewProductionEncoderConfig()
	err = l.CustomEncoderConfig(&encoderConfig)
	return
}

// NewEncoderConfigTextDevDefault   开发环境默认的 EncoderConfig 配置
func NewEncoderConfigTextDevDefault(l *Logger) (encoderConfig zapcore.EncoderConfig, err error) {
	encoderConfig = zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	return
}

func NewEncoderConfigTextDevCustom(l *Logger) (encoderConfig zapcore.EncoderConfig, err error) {
	encoderConfig = zap.NewDevelopmentEncoderConfig()
	err = l.CustomEncoderConfig(&encoderConfig)
	return
}

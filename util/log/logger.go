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

// NewEncoderConfigTextDevDefault   生产环境默认的 EncoderConfig 配置
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

type Logger struct {
	logConfig *LoggerConfig // 配置
	Logger    *zap.Logger   // 记录器实例
}

// NewLogger
func NewLogger(loggerConfig *LoggerConfig) (log *Logger, err error) {

	log = &Logger{
		logConfig: loggerConfig,
		Logger:    nil,
	}
	if err = log.checkLoggerConfig(); err != nil {
		return
	}
	err = log.NewZapLogger()
	return
}

// checkLoggerConfig 检查 要生成新的 日志记录器的配置
func (l *Logger) checkLoggerConfig() (err error) {
	if l.logConfig == nil {
		err = fmt.Errorf("LoggerConfig cannot be nil")
		return
	}
	// 检查 LowLevelFile
	if l.logConfig.LowLevelFile == nil {
		err = fmt.Errorf("LowLevelFile cannot be nil")
		return
	}
	if err = l.logConfig.LowLevelFile.checkFileName(); err != nil {
		return
	}
	// 检查 HighLevelFile 是否需要设置 LowLevelFile; FileName 不会一样
	if l.logConfig.HighLevelFile == nil {
		l.logConfig.createHighLevelToLowLevel()
	}

	// 部分字符串类型的字段进行去除空格+转小写的操作
	l.logConfig.simpleFormat()

	// 检查 编码器 重要配置
	if err = l.logConfig.checkEncoderConfigText(); err != nil {
		return
	}
	if err = l.logConfig.checkEncoderText(); err != nil {
		return
	}

	return
}

// NewZapLogger 构建 zap Logger 实例
func (l *Logger) NewZapLogger() (err error) {
	var (
		encoderConfig zapcore.EncoderConfig
		encoder       zapcore.Encoder
	)
	// 获取  zapcore.EncoderConfig
	if encoderConfig, err = l.GetEncoderConfig(); err != nil {
		return
	}
	// 获取 zapcore.Encoder
	if encoder, err = l.GetEncoder(encoderConfig); err != nil {
		return
	}
	if l.logConfig.OutputConsole {
	}
	if l.logConfig.OutputFile {

	}

	if l.logConfig.SplitWriteFromLevel {

	}

	// TODO 待完成 zapcore.NewCore 和 zapcorere,.NewTee
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

// customEncoderConfig 配置自动定义 EncoderConfig
func (l *Logger) CustomEncoderConfig(encoderConfig *zapcore.EncoderConfig) (err error) {
	// 配置各种 编码器
	var EncodeTime zapcore.TimeEncoder
	if err = EncodeTime.UnmarshalText([]byte(l.logConfig.TimeEncoderText)); err != nil {
		return
	}
	encoderConfig.EncodeTime = EncodeTime

	var EncodeLevel zapcore.LevelEncoder
	if err = EncodeLevel.UnmarshalText([]byte(l.logConfig.LevelEncoderText)); err != nil {
		return
	}
	encoderConfig.EncodeLevel = EncodeLevel

	var EncodeDuration zapcore.DurationEncoder
	if err = EncodeDuration.UnmarshalText([]byte(l.logConfig.DurationEncoderText)); err != nil {
		return
	}
	encoderConfig.EncodeDuration = EncodeDuration

	var EncodeCaller zapcore.CallerEncoder
	if err = EncodeCaller.UnmarshalText([]byte(l.logConfig.CallerEncoderText)); err != nil {
		return
	}
	encoderConfig.EncodeCaller = EncodeCaller

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

// 2个文件+1个控制台输出

/*
func getLogger(infoPath,errorPath string)  (*zap.Logger,error) {
	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool{
		return lev >= zap.ErrorLevel
	})

	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev < zap.ErrorLevel && lev >= zap.DebugLevel
	})


	prodEncoder := zap.NewProductionEncoderConfig()
	prodEncoder.EncodeTime = zapcore.ISO8601TimeEncoder


	lowWriteSyncer,lowClose,err :=  zap.Open(infoPath)
	if err != nil {
		lowClose()
		return nil,err
	}


	highWriteSyncer,highClose,err :=  zap.Open(errorPath)
	if err != nil {
		highClose()
		return nil,err
	}

	highCore := zapcore.NewCore(zapcore.NewJSONEncoder(prodEncoder),highWriteSyncer,highPriority)
	lowCore := zapcore.NewCore(zapcore.NewJSONEncoder(prodEncoder),lowWriteSyncer,lowPriority)


	return  zap.New(zapcorere,.NewTee(highColowCore),zap.AddCaller()),nil

*/

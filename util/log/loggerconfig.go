/**
@Author: wei-g
@Date:   2020/6/29 11:27 上午
@Description:
*/

package log

import (
	"fmt"
	"github.com/twgcode/sparrow/util/sliceutil"
	"path/filepath"
	"strings"
)

// 先看下 zap.AddStacktrace，研究完毕
//

/*
	TODO 研究下 zap.AddStacktrace
	已经研究明白
	logger, _ := zap.NewProduction(zap.AddStacktrace(zapcore.WarnLevel))
*/

/*
func NewDevelopmentEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
]
func NewProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
*/

const (
	// 默认的日志名
	DefaultHighLevelFileName = "error.log"
	DefaultLowLevelFileName  = "info.log"
)
const (
	// 选择 LoggerConfig.EncoderText 时使用
	EncoderTextJSON    = "json"
	EncoderTextConsole = "console"

	// 选择环境 LoggerConfig.EncoderConfigText
	EncoderConfigTextProdCustom  = "prod"
	EncoderConfigTextProdDefault = "prod_default"
	EncoderConfigTextDevCustom   = "dev"
	EncoderConfigTextDevDefault  = "dev_default"
)

var (
	// LoggerConfig.EncoderConfigText 可选值
	encoderConfigTextList = []string{EncoderConfigTextProdCustom, EncoderConfigTextProdDefault, EncoderConfigTextDevCustom, EncoderConfigTextDevDefault}
	encoderTextList       = []string{EncoderTextJSON, EncoderTextConsole}
)

func optionalValuesErr(value, class string) (err error) {
	err = fmt.Errorf("%s not a legal value, please check the optional value of %s", value, class)
	return
}

// LoggerConfig  提供生成一个 zap.Logger + lumberjack 的配置
type LoggerConfig struct {
	OutputConsole       bool // 是否输出到控制台
	OutputFile          bool // 是否输出到日志文件
	SplitWriteFromLevel bool // 是否根据不同的日志级别写不同的日志

	LowLevel                    string // 低级别日志等级, 不同级别的日志 写入的不同日志文件,
	HighLevel                   string // 同时高级别的日志也会记录 堆栈信息
	LowLevelFile, HighLevelFile *LoggerFileConfig

	TimeEncoderText     string // 时间格式编码器, 输出时间的格式
	LevelEncoderText    string // 日志等级编码器; A LevelEncoder serializes a Level to a primitive type.
	DurationEncoderText string // 持续时间编码器; A DurationEncoder serializes a time.Duration to a primitive type.
	CallerEncoderText   string // 调用文本编码器; 输出文件信息时，是以/full/path/to/package/file:line 全路径还是 package/file:line 的短路径

	TimeKey       string // 输出日志时时间的key
	LevelKey      string // 日志级别的key
	CallerKey     string // 调用文本的key; file:line
	MessageKey    string // 日志内容的key
	StacktraceKey string // 堆栈信息的key

	EncoderText       string // 日志编码器, 用来决定日志记录的整体形式; 有 json 和 console 2 种
	EncoderConfigText string // 日志配置编码器, An EncoderConfig allows users to configure the concrete encoders supplied by zapcore.
}

// simpleFormat 部分字符串类型的字段进行去除空格+转小写的操作,为了让 项目/业务代码更快 run 起来 没使用反射，如果后期压测影响不太这快可用反射减少代码量
func (l *LoggerConfig) simpleFormat() {
	changeField := []string{l.LowLevel, l.HighLevel, l.TimeEncoderText, l.LevelEncoderText, l.DurationEncoderText, l.CallerEncoderText, l.TimeKey, l.LevelKey, l.CallerKey, l.MessageKey, l.EncoderText, l.EncoderConfigText}
	for i := range changeField {
		changeField[i] = strings.ToLower(strings.TrimSpace(changeField[i]))
	}
}

// createHighLevelToLowLevel 根据 LowLevelFile 创建 HighLevelFile
func (l *LoggerConfig) createHighLevelToLowLevel() {
	if l.LowLevelFile == nil {
		return
	}
	_tmp := *l.LowLevelFile
	dir, _ := filepath.Split(_tmp.FileName)
	// 设置个默认的错误日志名
	_tmp.FileName = filepath.Join(dir, DefaultHighLevelFileName)
	l.HighLevelFile = &_tmp
}

// checkEncoderConfigText 检查和格式化 EncoderConfigText
func (l *LoggerConfig) checkEncoderConfigText() (err error) {
	// 不存在
	if !sliceutil.ContainsElement(encoderConfigTextList, l.EncoderConfigText) {
		err = optionalValuesErr(l.EncoderConfigText, "EncoderConfigText")
		return
	}
	return
}

func (l *LoggerConfig) checkEncoderText() (err error) {
	// 不存在
	if !sliceutil.ContainsElement(encoderTextList, l.EncoderText) {
		err = optionalValuesErr(l.EncoderText, "EncoderText")
		return
	}
	return
}

// LoggerFile 日志文件
type LoggerFileConfig struct {
	FileName                    string
	MaxSize, MaxBackups, MaxAge int
	Compress                    bool // 是否开启日志压缩，默认为否
}

func (l LoggerFileConfig) checkFileName() (err error) {
	l.FileName = strings.TrimSpace(l.FileName)
	if l.FileName == "" {
		return fmt.Errorf("FileName Cannot be empty")
	}

	return
}

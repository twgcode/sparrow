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
	"sync"
)

const (
	// 默认的日志名
	DefaultHighLevelFileName = "error"
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
	OutputConsole       bool `mapstructure:"output_console"`         // 是否输出到控制台
	OutputFile          bool `mapstructure:"output_file"`            // 是否输出到日志文件
	SplitWriteFromLevel bool `mapstructure:"split_write_from_level"` // 是否根据不同的日志级别写不同的日志
	Stacktrace          bool `mapstructure:"stacktrace"`             // 是否 记录高级别日志时记录对应的堆栈信息

	LowLevel      string            `mapstructure:"low_level"`  // 低级别日志等级, 不同级别的日志 写入的不同日志文件,
	HighLevel     string            `mapstructure:"high_level"` // 同时高级别的日志也会记录 堆栈信息
	LowLevelFile  *LoggerFileConfig `mapstructure:"low_level_file"`
	HighLevelFile *LoggerFileConfig `mapstructure:"high_level_file"`

	TimeEncoderText     string `mapstructure:"time_encoder_text"`     // 时间格式编码器, 输出时间的格式
	LevelEncoderText    string `mapstructure:"level_encoder_text"`    // 日志等级编码器; A LevelEncoder serializes a Level to a primitive type.
	DurationEncoderText string `mapstructure:"duration_encoder_text"` // 持续时间编码器; A DurationEncoder serializes a time.Duration to a primitive type.
	CallerEncoderText   string `mapstructure:"caller_encoder_text"`   // 调用文本编码器; 输出文件信息时，是以/full/path/to/package/file:line 全路径还是 package/file:line 的短路径

	TimeKey       string `mapstructure:"time_key"`       // 输出日志时时间的key
	LevelKey      string `mapstructure:"level_key"`      // 日志级别的key
	CallerKey     string `mapstructure:"callel_key"`     // 调用文本的key; file:line
	MessageKey    string `mapstructure:"message_key"`    // 日志内容的key
	StacktraceKey string `mapstructure:"stacktrace_key"` // 堆栈信息的key

	CallerSkip int `mapstructure:"caller_skip"` // 跳过几层调用者,在对zap输出的方法进行封装时会使用该参数

	EncoderText       string `mapstructure:"encoder_text"`        // 日志编码器, 用来决定日志记录的整体形式; 有 json 和 console 2 种
	EncoderConfigText string `mapstructure:"encoder_config_text"` // 日志配置编码器, An EncoderConfig allows users to configure the concrete encoders supplied by zapcore.
	checkMutex        sync.Mutex
}

// checkLoggerConfig 检查 要生成新的 日志记录器的配置
func (l *LoggerConfig) CheckLoggerConfig() (err error) {
	defer l.checkMutex.Unlock()
	l.checkMutex.Lock()
	// 检查 LowLevelFile
	if !l.OutputFile && !l.OutputConsole {
		err = fmt.Errorf("please set the correct output source, set one of LoggerConfig.OutputConsole or LoggerConfig.OutputFile")
		return
	}
	if l.OutputFile {
		if l.LowLevelFile == nil {
			err = fmt.Errorf("LowLevelFile cannot be nil")
			return
		}
		if err = l.LowLevelFile.checkFileName(); err != nil {
			return
		}
		// 检查 HighLevelFile 是否需要设置 LowLevelFile; FileName 不会一样
		if l.HighLevelFile == nil && l.SplitWriteFromLevel {
			l.createHighLevelToLowLevel()
		}
	}

	// 部分字符串类型的字段进行去除空格+转小写的操作
	l.simpleFormat()

	// 检查 编码器 重要配置
	if err = l.checkEncoderConfigText(); err != nil {
		return
	}
	if err = l.checkEncoderText(); err != nil {
		return
	}
	// 检查 配置编码器是时需要的各种keys
	err = l.checkKeys()
	return
}

// simpleFormat 部分字符串类型的字段进行去除空格+转小写的操作,为了让 项目/业务代码更快 run 起来 没使用反射，如果后期压测影响不太这快可用反射减少代码量
func (l *LoggerConfig) simpleFormat() {
	if l.EncoderConfigText == EncoderConfigTextProdCustom || l.EncoderConfigText == EncoderConfigTextDevCustom {
		changeField := []*string{&l.LowLevel, &l.HighLevel, &l.TimeEncoderText, &l.LevelEncoderText, &l.DurationEncoderText, &l.CallerEncoderText, &l.EncoderText, &l.EncoderConfigText}
		for i := range changeField {
			*changeField[i] = strings.ToLower(strings.TrimSpace(*changeField[i]))
		}
		// 只进行去除字符串
		changeField = []*string{&l.TimeKey, &l.LevelKey, &l.CallerKey, &l.MessageKey, &l.StacktraceKey}
		for i := range changeField {
			*changeField[i] = strings.TrimSpace(*changeField[i])
		}
	}

}

// checkKeys 检查配置 zapcore.EncoderConfig时的 各种key,注意要在执行checkEncoderConfigText后在执行改方法
func (l *LoggerConfig) checkKeys() (err error) {
	if l.EncoderConfigText == EncoderConfigTextProdCustom || l.EncoderConfigText == EncoderConfigTextDevCustom {

		for _, v := range []string{l.TimeKey, l.LevelKey, l.CallerKey, l.MessageKey, l.StacktraceKey} {
			if len(v) == 0 {
				err = fmt.Errorf("an illegal value to configure zapcore.EncoderConfig Key key cannot be empty, value is %s", v)
				return
			}
		}
	}
	return
}

// createHighLevelToLowLevel 根据 LowLevelFile 创建 HighLevelFile
func (l *LoggerConfig) createHighLevelToLowLevel() {
	if l.LowLevelFile == nil {
		return
	}
	_tmp := *l.LowLevelFile
	dir, fileName := filepath.Split(_tmp.FileName)

	fileNameList := strings.Split(fileName, ".")
	// 设置个默认的错误日志名
	if len(fileNameList) == 1 {
		_tmp.FileName = filepath.Join(dir, fileName+"-"+DefaultHighLevelFileName)
	} else {
		name := fmt.Sprintf("%s-%s.%s", strings.Join(fileNameList[0:len(fileNameList)-1], "."), DefaultHighLevelFileName, fileNameList[len(fileNameList)-1])
		_tmp.FileName = filepath.Join(dir, name)
	}
	l.HighLevelFile = &_tmp
}

// checkEncoderConfigText 检查 EncoderConfigText
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
	FileName   string `mapstructure:"file_name"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"` // 是否开启日志压缩，默认为否
}

func (l LoggerFileConfig) checkFileName() (err error) {
	l.FileName = strings.TrimSpace(l.FileName)
	if l.FileName == "" {
		return fmt.Errorf("FileName Cannot be empty")
	}

	return
}

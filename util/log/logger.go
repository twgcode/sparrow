/*
@Author: wei-g
@Email:  wei_g_it@163.com
@Date:   2020/6/28 下午9:30
@Description:
*/

package log

/*
	TODO 研究下 zap.AddStacktrace
	logger, _ := zap.NewProduction(zap.AddStacktrace(zapcore.WarnLevel))
*/

/*
{"level":"INFO","time":"2020-04-28T21:08:58.504+0800","caller":"middleware/logger.go:25","msg":"/down_data/v1/pc/deploy","status":200,"method":"GET","path":"/down_data/v1/pc/deploy","query":"page_num=1&page_size=10","ip":"223.71.139.100","user-agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.122 Safari/537.36","errors":"","cost":0.023618231}
*/

/*
	TODO. 下面的这些key还没整理到 Logger 结构体中
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "caller",
	MessageKey:     "msg",
	StacktraceKey:  "stacktrace",

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

// Logger 提供日志记录服务; Logger 融合了 zap 和 lumberjack
type Logger struct {
	OutputConsole       bool // 是否输出到控制
	OutputFile          bool // 是否输出到日志文件
	SplitWriteFromLevel bool // 是否根据不同的日志级别写不同的日志

	LowLevelFile, HighLevelFile *LoggerFile // 不同级别的日志 写入的不同日志文件

	TimeEncoderText     string // 时间格式编码器, 输出时间的格式
	LevelEncoderText    string // 日志等级编码器; A LevelEncoder serializes a Level to a primitive type.
	DurationEncoderText string // 持续时间编码器; A DurationEncoder serializes a time.Duration to a primitive type.
	CallerEncoderText   string // 调用文本编码器; 输出文件信息时，是以/full/path/to/package/file:line 全路径还是 package/file:line 的短路径

	TimeKey string // 输出日志时时间的key
}

// LoggerFile 日志文件
type LoggerFile struct {
	FileName                    string
	MaxSize, MaxBackups, MaxAge int
	Compress                    bool // 是否开启日志压缩，默认为否
	//
}

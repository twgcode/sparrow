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
	if encoderConfig, err = l.logConfig.getEncoderConfig(); err != nil {
		return
	}
	// 获取 zapcore.Encoder
	if encoder, err = l.logConfig.getEncoder(encoderConfig); err != nil {
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

// getFileLogWriter 结合日志切割工具设置日志输出接口
func (l *Logger) getFileLogWriter(fileConfig *LoggerFileConfig) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fileConfig.FileName,
		MaxSize:    fileConfig.MaxSize,
		MaxBackups: fileConfig.MaxBackups,
		MaxAge:     fileConfig.MaxAge,
		Compress:   fileConfig.Compress,
	}
	return zapcore.AddSync(lumberJackLogger)
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

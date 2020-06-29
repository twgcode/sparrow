/**
@Author: wei-g
@Date:   2020/6/29 5:51 下午
@Description:
*/

package log

import "testing"

func TestSimpleFormat(t *testing.T) {
	logger := &LoggerConfig{
		OutputConsole:       false,
		OutputFile:          false,
		SplitWriteFromLevel: false,
		LowLevel:            "LowLevel",
		HighLevel:           "HighLevel",
		LowLevelFile:        nil,
		HighLevelFile:       nil,
		TimeEncoderText:     "TimeEncoderText",
		LevelEncoderText:    "LevelEncoderText",
		DurationEncoderText: "DurationEncoderText",
		CallerEncoderText:   "CallerEncoderText",
		TimeKey:             "TimeKey",
		LevelKey:            "LevelKey",
		CallerKey:           "CallerKey",
		MessageKey:          "MessageKey",
		StacktraceKey:       "StacktraceKey",
		EncoderText:         "EncoderText",
		EncoderConfigText:   "EncoderConfigText",
	}
	logger.simpleFormat()
	t.Logf("%#v\n", logger)
}

package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func GetLogger() *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 使用彩色日志等级
		EncodeTime:     zapcore.ISO8601TimeEncoder,       // 设置时间格式
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	// 创建一个 console 输出的 Core
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig), // 使用 ConsoleEncoder 进行输出
		zapcore.AddSync(zapcore.Lock(os.Stdout)), // 输出到标准输出
		zapcore.DebugLevel,                       // 设置日志级别
	)

	// 创建一个 Logger
	return zap.New(core, zap.AddCaller())
}

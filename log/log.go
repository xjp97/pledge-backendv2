package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
	"runtime"
)

var Logger *zap.Logger

func init() {
	// zap 不支持文件归档, 如果要支持文件按大小或者时间归档, 需要使用 lumberjack
	// https://github.com/uber-go/zap/blob/master/FAQ.md
	hook := lumberjack.Logger{
		Filename:   getCurrentAbPathByCaller() + "log/pledge.log", // 日志存储目录
		MaxSize:    500,                                           //每个文件保存的最大尺寸,单位M
		MaxBackups: 20,                                            // 日志文件最多保存多少个备份
		MaxAge:     7,                                             // 文件最多保存多少天
		Compress:   true,                                          // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // IS08601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, // 编码时长
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,        // 编码名称
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.InfoLevel)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)

	// 开启开发模式,栈堆跟踪
	caller := zap.AddCaller()
	// 开启文件以及行号
	development := zap.Development()
	// 设置初始化字段
	field := zap.Fields(zap.String("serviceName", "pledge"))
	// 构造日志
	Logger = zap.New(core, caller, development, field)
}

// 获取当前执行文件绝对路径（go run）
func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

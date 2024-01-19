package log

import (
	"ai-report/common"
	"ai-report/common/consts"
	"ai-report/common/utils"
	"ai-report/entity"
	"context"
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

func InitZap() {
	// 判断是否有Director文件夹
	if ok := utils.PathExists(common.GlobalConfig.Zap.Director); !ok {
		fmt.Printf("create %v directory\n", common.GlobalConfig.Zap.Director)
		_ = os.Mkdir(common.GlobalConfig.Zap.Director, os.ModePerm)
	}
	// 调试级别
	debugPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.DebugLevel
	})
	// 日志级别
	infoPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.InfoLevel
	})
	// 警告级别
	warnPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.WarnLevel
	})
	// 错误级别
	errorPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.ErrorLevel
	})
	cores := [...]zapcore.Core{
		getEncoderCore(fmt.Sprintf("./%s/server_debug.log", common.GlobalConfig.Zap.Director), debugPriority),
		getEncoderCore(fmt.Sprintf("./%s/server_info.log", common.GlobalConfig.Zap.Director), infoPriority),
		getEncoderCore(fmt.Sprintf("./%s/server_warn.log", common.GlobalConfig.Zap.Director), warnPriority),
		getEncoderCore(fmt.Sprintf("./%s/server_error.log", common.GlobalConfig.Zap.Director), errorPriority),
	}
	logger := zap.New(zapcore.NewTee(cores[:]...), zap.AddCaller(), zap.AddCallerSkip(1))
	common.Logger = logger
	Log = LogWrapper{logger: logger}
}

// getEncoderConfig 获取zapcore.EncoderConfig
func getEncoderConfig() (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  common.GlobalConfig.Zap.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	switch {
	case common.GlobalConfig.Zap.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case common.GlobalConfig.Zap.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case common.GlobalConfig.Zap.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case common.GlobalConfig.Zap.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return config
}

// getEncoder 获取zapcore.Encoder
func getEncoder() zapcore.Encoder {
	if common.GlobalConfig.Zap.Format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}

// getEncoderCore 获取Encoder的zapcore.Core
func getEncoderCore(fileName string, level zapcore.LevelEnabler) (core zapcore.Core) {
	writer := GetWriteSyncer(fileName) // 使用file-rotatelogs进行日志分割
	return zapcore.NewCore(getEncoder(), writer, level)
}

// 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006/01/02 15:04:05"))
}

// @function: GetWriteSyncer
// @description: zap logger中加入file-rotatelogs
// @return: zapcore.WriteSyncer, error
func GetWriteSyncer(file string) zapcore.WriteSyncer {
	// 每小时一个文件
	logf := lumberjack.Logger{
		Filename:   "logs/log.log", // 日志文件路径
		MaxSize:    128,            // 每个日志文件保存的大小 单位:M
		MaxAge:     60,             // 文件最多保存多少天
		MaxBackups: 30,             // 日志文件最多保存多少个备份
		Compress:   false,          // 是否压缩
	}
	if common.GlobalConfig.Zap.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&logf))
	}
	return zapcore.AddSync(&logf)
}

type LogWrapper struct {
	logger *zap.Logger
}

var Log LogWrapper

func Debug(tag string, fields ...zap.Field) {
	Log.logger.Debug(tag, fields...)
}
func DebugF(ctx context.Context, tag string, fields ...zap.Field) {
	trace := ctx.Value(consts.TraceKey).(*entity.Trace)
	Log.logger.Debug(tag,
		append(fields, zap.String("trace_id", trace.TraceId))...,
	)
}
func Info(tag string, fields ...zap.Field) {
	Log.logger.Info(tag, fields...)
}
func InfoF(ctx context.Context, tag string, fields ...zap.Field) {
	trace := ctx.Value(consts.TraceKey).(*entity.Trace)
	Log.logger.Info(tag,
		append(fields, zap.String("trace_id", trace.TraceId))...,
	)
}
func Warn(tag string, fields ...zap.Field) {
	Log.logger.Warn(tag, fields...)
}
func WarnF(ctx context.Context, tag string, fields ...zap.Field) {
	trace := ctx.Value(consts.TraceKey).(*entity.Trace)
	Log.logger.Warn(tag,
		append(fields, zap.String("trace_id", trace.TraceId))...,
	)
}
func Error(tag string, fields ...zap.Field) {
	Log.logger.Error(tag, fields...)
}
func ErrorF(ctx context.Context, tag string, fields ...zap.Field) {
	trace := ctx.Value(consts.TraceKey).(*entity.Trace)
	Log.logger.Error(tag,
		append(fields, zap.String("trace_id", trace.TraceId))...,
	)
}
func Fatal(tag string, fields ...zap.Field) {
	Log.logger.Fatal(tag, fields...)
}
func FatalF(ctx context.Context, tag string, fields ...zap.Field) {
	trace := ctx.Value(consts.TraceKey).(*entity.Trace)
	Log.logger.Fatal(tag,
		append(fields, zap.String("trace_id", trace.TraceId))...,
	)
}

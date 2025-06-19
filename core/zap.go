package core

import (
	"fmt"
	"git.exclouds.org/ww/greet_ws/global"
	"git.exclouds.org/ww/greet_ws/utils"
	zaprotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var directoryUtils utils.DirectoryUtils

var (
	err    error
	level  zapcore.Level
	writer zapcore.WriteSyncer
)

func init() {
	if ok, _ := directoryUtils.PathExists(global.SERVER_CONFIG.ZapConfig.Directory); !ok { // 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", global.SERVER_CONFIG.ZapConfig.Directory)
		_ = os.Mkdir(global.SERVER_CONFIG.ZapConfig.Directory, os.ModePerm)
	}

	switch global.SERVER_CONFIG.ZapConfig.Level { // 初始化配置文件的Level
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}

	writer, err = getWriteSyncer() // 使用file-rotatelogs进行日志分割
	if err != nil {
		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
		return
	}

	if level == zap.DebugLevel || level == zap.ErrorLevel {
		global.LOGGER = zap.New(getEncoderCore(), zap.AddStacktrace(level))
	} else {
		global.LOGGER = zap.New(getEncoderCore())
	}
	if global.SERVER_CONFIG.ZapConfig.ShowLine {
		global.LOGGER.WithOptions(zap.AddCaller())
	}
}

// getWriteSyncer zap logger中加入file-rotatelogs
func getWriteSyncer() (zapcore.WriteSyncer, error) {
	zapConfig := global.SERVER_CONFIG.ZapConfig
	fileWriter, err := zaprotatelogs.New(
		zapConfig.Directory+string(os.PathSeparator)+"%Y-%m-%d-%H.log",
		//zaprotatelogs.WithLinkName(global.SERVER_CONFIG.ZapConfig.LinkName),//create soft link
		zaprotatelogs.WithMaxAge(zapConfig.KeepDays*24*time.Hour),
		zaprotatelogs.WithRotationTime(time.Hour),
	)
	if zapConfig.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)), err
	}
	return zapcore.AddSync(fileWriter), err
}

// getEncoderConfig 获取zapcore.EncoderConfig
func getEncoderConfig() (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  global.SERVER_CONFIG.ZapConfig.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	switch {
	case global.SERVER_CONFIG.ZapConfig.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case global.SERVER_CONFIG.ZapConfig.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case global.SERVER_CONFIG.ZapConfig.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case global.SERVER_CONFIG.ZapConfig.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	return config
}

// getEncoder 获取zapcore.Encoder
func getEncoder() zapcore.Encoder {
	if global.SERVER_CONFIG.ZapConfig.Format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}

// getEncoderCore 获取Encoder的zapcore.Core
func getEncoderCore() (core zapcore.Core) {
	return zapcore.NewCore(getEncoder(), writer, level)
}

// 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(global.SERVER_CONFIG.ZapConfig.Prefix + "2006/01/02 - 15:04:05.000"))
}

package config

import "time"

type ZapConfig struct {
	Level         string        `mapstructure:"level" json:"level" yaml:"level"`
	Format        string        `mapstructure:"format" json:"format" yaml:"format"`
	Prefix        string        `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	Directory     string        `mapstructure:"directory" json:"directory"  yaml:"directory"`
	LinkName      string        `mapstructure:"link-name" json:"linkName" yaml:"link-name"`
	ShowLine      bool          `mapstructure:"show-line" json:"showLine" yaml:"showLine"`
	EncodeLevel   string        `mapstructure:"encode-level" json:"encodeLevel" yaml:"encode-level"`
	StacktraceKey string        `mapstructure:"stacktrace-key" json:"stacktraceKey" yaml:"stacktrace-key"`
	LogInConsole  bool          `mapstructure:"log-in-console" json:"logInConsole" yaml:"log-in-console"`
	KeepDays      time.Duration `mapstructure:"keep-days" json:"keepDays" yaml:"keep-days"`             //日志保留天数
	RotationTime  time.Duration `mapstructure:"rotation-time" json:"rotationTime" yaml:"rotation-time"` // 新增：日志轮转时间
}

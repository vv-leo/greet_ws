package config

type ServerConfig struct {
	SystemConfig SystemConfig `mapstructure:"system" json:"system" yaml:"system"`
	RedisConfig  RedisConfig  `mapstructure:"redis" json:"redis" yaml:"redis"`
	ZapConfig    ZapConfig    `mapstructure:"zap" json:"zap" yaml:"zap"`
	// consul
	ConsulConfig ConsulConfig `mapstructure:"consul" json:"consul" yaml:"consul"`
}

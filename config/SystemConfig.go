package config

type SystemConfig struct {
	Env               string   `mapstructure:"env" json:"env" yaml:"env"`
	Addr              int      `mapstructure:"addr" json:"addr" yaml:"addr"`
	DbType            string   `mapstructure:"db-type" json:"db-type" yaml:"db-type"`
	UseMultipoint     bool     `mapstructure:"use-multipoint" json:"use-multipoint" yaml:"use-multipoint"`
	AllowOrigins      []string `mapstructure:"allow-origins" json:"allow-origins" yaml:"allow-origins"`
	MaxKeepTime       int64    `mapstructure:"max-keep-time" json:"max-keep-time" yaml:"max-keep-time"`
	ShardedTableCount int64    `mapstructure:"sharded-table-count" json:"sharded-table-count" yaml:"sharded-table-count"`
	Token             string   `mapstructure:"token" json:"token" yaml:"token"`
}

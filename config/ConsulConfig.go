package config

type ConsulConfig struct {
	Address    string      `mapstructure:"address" json:"address" yaml:"address"`
	Datacenter string      `mapstructure:"datacenter" json:"datacenter" yaml:"datacenter"`
	Service    ServiceInfo `mapstructure:"service" json:"service" yaml:"service"`
}

type ServiceInfo struct {
	ID      string       `mapstructure:"id" json:"id" yaml:"id"`
	Name    string       `mapstructure:"name" json:"name" yaml:"name"`
	Tags    []string     `mapstructure:"tags" json:"tags" yaml:"tags"`
	Check   ServiceCheck `mapstructure:"check" json:"check" yaml:"check"`
	Address string       `mapstructure:"address" json:"address" yaml:"address"`
}

type ServiceCheck struct {
	HTTP                           string `mapstructure:"http" json:"http" yaml:"http"`
	Interval                       string `mapstructure:"interval" json:"interval" yaml:"interval"`
	Timeout                        string `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
	DeregisterCriticalServiceAfter string `mapstructure:"deregister-critical-service-after" json:"deregister-critical-service-after" yaml:"deregister-critical-service-after"`
}

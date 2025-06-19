package config

type EndpointsConfig struct {
	WsServer WsServer `mapstructure:"ws-server" json:"ws-server" yaml:"ws-server"`
}

type WsServer struct {
	HttpUrl      string `mapstructure:"http-url" json:"http-url" yaml:"http-url"`
	WebsocketUrl string `mapstructure:"websocket-url" json:"websocket-url" yaml:"websocket-url"`
	IntervalTime int    `mapstructure:"interval-time" json:"interval-time" yaml:"interval-time"`
}

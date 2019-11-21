package common

// AppConfig app configurations
type AppConfig struct {
	LogCfg *LogConfig
	SrvCfg *ServerConfig
}

// ServerConfig syslog server configuration
type ServerConfig struct {
	TCPAddr string // listen ip
	UDPAddr string
	RFC3164 string // port num for syslog format RFC3164
	RFC5424 string //
	Filters []string
}

// AppConfigInstance ...
var AppConfigInstance *AppConfig

func init() {
	AppConfigInstance = &AppConfig{
		LogCfg: &LogConfig{},
		SrvCfg: &ServerConfig{Filters: make([]string, 0)},
	}
}

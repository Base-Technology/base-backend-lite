package conf

var Conf *Config

type Config struct {
	ServerConf       ServerConfig       `mapstructure:"server"`
	DBConf           DBConfig           `mapstructure:"database"`
	LoggerConf       LoggerConfig       `mapstructure:"logger"`
	ChatGPTProxyConf ChatGPTProxyConfig `mapstructure:"chatgpt_porxy"`
}

type ServerConfig struct {
	Port            int    `mapstructure:"port"`
	TokenExpireTime int    `mapstructure:"token_expire_time"`
	TokenSecret     string `mapstructure:"token_secret"`
}

type DBConfig struct {
	IP       string `mapstructure:"ip"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

type ChatGPTProxyConfig struct {
	IP   string `mapstructure:"ip"`
	Port int    `mapstructure:"port"`
}

const (
	BaseBackendLiteConfigPrefix      = "BASE_BACKEND_LITE"
	DefaultBaseBackendLiteConfigFile = "config/base_backend_lite_config.yaml"
)

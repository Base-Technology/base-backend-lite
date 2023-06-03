package conf

var Conf *Config

type Config struct {
	ServerConf       ServerConfig       `mapstructure:"server"`
	DBConf           DBConfig           `mapstructure:"database"`
	LoggerConf       LoggerConfig       `mapstructure:"logger"`
	ChatGPTProxyConf ChatGPTProxyConfig `mapstructure:"chatgpt_porxy"`
	ValidateCodeConf ValidateCodeConfig `mapstructure:"validate_code"`
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

type ValidateCodeConfig struct {
	Server    string `mapstructure:"server"`
	AppKey    string `mapstructure:"appkey"`
	AppSecret string `mapstructure:"appsecret"`
	AppCode   string `mapstructure:"appcode"`
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

package conf

var Conf *Config

type Config struct {
	ServerConf       ServerConfig       `mapstructure:"server"`
	DBConf           DBConfig           `mapstructure:"database"`
	ValidateCodeConf ValidateCodeConfig `mapstructure:"validate_code"`
	IMTPConf         IMTPConfig         `mapstructure:"imtp"`
	LoggerConf       LoggerConfig       `mapstructure:"logger"`
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

type IMTPConfig struct {
	APPServer string `mapstructure:"app_server"`
	APIServer string `mapstructure:"api_server"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

const (
	BaseBackendLiteConfigPrefix      = "BASE_BACKEND_LITE"
	DefaultBaseBackendLiteConfigFile = "config/base_backend_lite_config.yaml"
)

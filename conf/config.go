package conf

import (
	"strings"

	"github.com/Base-Technology/base-app-lite/seelog"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// InitConfig 初始化配置
func InitConfig() error {
	v := viper.New()
	v.SetEnvPrefix(BaseBackendLiteConfigPrefix)
	v.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)
	v.SetConfigFile(DefaultBaseBackendLiteConfigFile)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		return errors.Errorf("read config error, %v", err)
	}
	for _, key := range v.AllKeys() {
		seelog.Infof("%s=%v", key, v.Get(key))
	}
	Conf = &Config{}
	if err := v.Unmarshal(Conf); err != nil {
		return errors.Errorf("unmarshal Config error, %v", err)
	}
	setLoggerLevel()
	return nil
}

func setLoggerLevel() {
	switch strings.ToUpper(Conf.LoggerConf.Level) {
	case "INFO":
		seelog.SetInfoLevel()
	case "WARN":
		seelog.SetWarnLevel()
	case "ERROR":
		seelog.SetErrorLevel()
	}
}

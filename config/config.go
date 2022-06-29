package config

import "os"

var (
	defaultPort = ":6969"
	defaultLng  = "ru-RU"
	defaultKey  = "6gHksNL9Zjck2zEx"
	defaultUrl  = "https://api.textgears.com"
)

type Config struct {
	Port string
	SpellerCfg
}

type SpellerCfg struct {
	Key      string
	Url      string
	Language string
}

func New() *Config {
	return &Config{}
}

func (c *Config) Load() {
	c.Port = setEnvStr("port", defaultPort)
	c.Key = setEnvStr("key", defaultKey)
	c.Language = setEnvStr("language", defaultLng)
	c.Url = setEnvStr("url", defaultUrl)
}

func setEnvStr(key, defaultStr string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultStr
}

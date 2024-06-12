package configs

import (
	"time"

	"github.com/spf13/viper"
)

type (
	configs struct {
		Port                      string           `mapstructure:"PORT"`
		Host                      string           `mapstructure:"DB_HOST"`
		IpLimit                   int              `mapstructure:"IP_LIMIT"`
		TokenLimit                int              `mapstructure:"TOKEN_LIMIT"`
		IpBlockTime               int              `mapstructure:"IP_BLOCK_TIME"`
		TokenBlockTime            int              `mapstructure:"TOKEN_BLOCK_TIME"`
		TokenMaxRequestsPerSecond map[string]int64 `mapstructure:"TOKEN_MAX_REQUESTS_PER_SECOND"`
	}

	LimitConfig struct {
		MaxRequestsIp             int
		MaxRequestsToken          int
		AllowedToken              string
		IpBlockTime               time.Duration
		TokenBlockTime            time.Duration
		TokenMaxRequestsPerSecond map[string]int64 `mapstructure:"TOKEN_MAX_REQUESTS_PER_SECOND"`
	}
)

func GetConfig(path string) (*configs, error) {
	var cfg *configs
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	return cfg, nil
}

func (c *configs) GetLimitConfig() LimitConfig {
	return LimitConfig{
		MaxRequestsIp:             c.IpLimit,
		MaxRequestsToken:          c.TokenLimit,
		TokenMaxRequestsPerSecond: c.TokenMaxRequestsPerSecond,
		IpBlockTime:               time.Duration(c.IpBlockTime) * time.Second,
		TokenBlockTime:            time.Duration(c.TokenBlockTime) * time.Second,
	}
}

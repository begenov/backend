package config

import (
	"time"

	"github.com/spf13/viper"
)

const (
	defaultServerPort               = "8080"
	defaultServerRWTimeout          = 10 * time.Second
	defaultServerMaxHeaderMegabytes = 1
	defaultAccessTokenDuration      = 15 * time.Minute
)

type Config struct {
	Postgres DBConfig
	Server   HTTPConfig
	JWT      JWTConfig
}

type DBConfig struct {
	Driver string `mapstructure:"DB_DRIVER"`
	DSN    string `mapstructure:"DB_SOURCE"`
}

type HTTPConfig struct {
	Addr           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
}

type JWTConfig struct {
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func Init(path string) (*Config, error) {
	viper.AddConfigPath(path)

	viper.SetConfigName("app")
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	viper.AutomaticEnv()

	cfg := Config{}
	if err := viper.UnmarshalKey("DB_DRIVER", &cfg.Postgres.Driver); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("DB_SOURCE", &cfg.Postgres.DSN); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("TOKEN_SYMMETRIC_KEY", &cfg.JWT.TokenSymmetricKey); err != nil {
		return nil, err
	}

	cfg.Server = HTTPConfig{
		Addr:           defaultServerPort,
		ReadTimeout:    defaultServerRWTimeout,
		WriteTimeout:   defaultServerRWTimeout,
		MaxHeaderBytes: defaultServerMaxHeaderMegabytes,
	}
	cfg.JWT.AccessTokenDuration = defaultAccessTokenDuration
	return &cfg, nil
}

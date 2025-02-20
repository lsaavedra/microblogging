package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

const (
	Local      string = "local"
	Develop    string = "dev"
	Production string = "prod"
)

type (
	Config struct {
		Environment   string
		ListeningPort string
		PgConfig
		CacheConfig
	}

	PgConfig struct {
		PgUser       string
		PgPassword   string
		PgDatabase   string
		PostgresPort string
		PostgresHost string
	}

	CacheConfig struct {
		Host string
		Port string
	}
)

func Load(env string) *Config {
	if env == "" {
		env = Local
	}

	viper.AutomaticEnv()
	viper.SetConfigName(fmt.Sprintf("config_%s", strings.ToLower(env)))
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./internal/config")
	viper.AddConfigPath("/app/internal/config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")
	viper.AddConfigPath("../internal/config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	viper.SetDefault("LISTENING_PORT", "5000")

	return &Config{
		Environment:   env,
		ListeningPort: viper.GetString("LISTENING_PORT"),
		PgConfig: PgConfig{
			PgUser:       viper.GetString("POSTGRES_USER"),
			PgPassword:   viper.GetString("POSTGRES_PASSWORD"),
			PgDatabase:   viper.GetString("POSTGRES_DB"),
			PostgresPort: viper.GetString("POSTGRES_PORT"),
			PostgresHost: viper.GetString("POSTGRES_HOST"),
		},
		CacheConfig: CacheConfig{
			Host: viper.GetString("REDIS_HOST"),
			Port: viper.GetString("REDIS_PORT"),
		},
	}
}

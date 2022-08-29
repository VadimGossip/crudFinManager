package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
	"log"
	"time"
)

type NetServerConfig struct {
	Host string
	Port int
}

type PostgresConfig struct {
	Host     string
	Port     int
	Username string
	Name     string
	SSLMode  string
	Password string
}

type AuthConfig struct {
	TokenTTL time.Duration `mapstructure:"token_ttl"`
	Salt     string
	Secret   string
}

type AuthConfig2 struct {
	Salt   string
	Secret string
}

type Config struct {
	Server   NetServerConfig
	Postgres PostgresConfig
	Auth     AuthConfig
}

func parseConfigFile(configDir string) error {
	viper.AddConfigPath(configDir)
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("serverListener.tcp", &cfg.Server); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("auth", &cfg.Auth); err != nil {
		return err
	}
	return nil
}

func setFromEnv(cfg *Config) error {
	if err := envconfig.Process("db", &cfg.Postgres); err != nil {
		return err
	}

	if err := envconfig.Process("auth", &cfg.Auth); err != nil {
		return err
	}
	return nil
}

func Init(configDir string) (*Config, error) {
	viper.SetConfigName("config")
	if err := parseConfigFile(configDir); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}
	if err := setFromEnv(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

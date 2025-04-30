package config

import (
	"errors"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env      string        `yaml:"env" env-default:"env"`
	TokenTTL time.Duration `yaml:"token_ttl" env-default:"1h"`
	GRPC     GRPCServer    `yaml:"grpc"`
	DB       DBConfig      `yaml:"db"`
}

type DBConfig struct {
	Port     string `yaml:"port"`
	Host     string `yaml:"host"`
	UserName string `yaml:"username"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
	Password string
}

type GRPCServer struct {
	Port    int           `yaml:"port" env-default:"8000"`
	Timeout time.Duration `yaml:"timeout" env-default:"10s"`
}

func InitConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err

	}

	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		return nil, errors.New("cannot read config")
	}

	if _, err := os.Stat(configPath); err != nil {
		return nil, errors.New("cannot read file")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil

}

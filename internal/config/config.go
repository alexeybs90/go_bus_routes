package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env     string  `yaml:"env" env-default:"local"`
	Server  Server  `yaml:"server"`
	Storage Storage `yaml:"storage"`
}

type Server struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"10s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Storage struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port" env-default:"5432"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
}

func LoadConfig(path string) (Config, error) {
	if path == "" {
		return Config{}, errors.New("config path is not set")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return Config{}, fmt.Errorf("config path is not exist: %s", path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return Config{}, fmt.Errorf("cannot read config: %s", path)
	}

	return cfg, nil
}

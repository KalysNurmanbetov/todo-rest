package config

import (
	"fmt"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `yaml:"env" env-required:"true"`
	HTTPServer
	Postgres
}

type HTTPServer struct {
	Address     string `yaml:"address" env-default:"localhost:9091"`
	Timeout     string `yaml:"timeout" env-default:"4s"`
	IdleTimeout string `yaml:"idle_timeout" env-default:"60s"`
}

type Postgres struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	DBname   string `yaml:"dbname" env-required:"true"`
	SSLMode  string `yaml:"sslmode" env-default:"disable"`
}

func (p Postgres) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", p.User, p.Password, p.Host, p.Port, p.DBname, p.SSLMode)
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatalf("enviroment variable CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var config Config
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("cannot read the config file: %s", err)
	}

	return &config
}

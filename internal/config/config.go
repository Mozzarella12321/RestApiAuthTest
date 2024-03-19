package config

import (
	"log"
	"os"
	"path"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {

	cfgPath := getCfgPath()

	var cfg Config

	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		log.Fatalf("Could not read config: %v", err)
	}

	return &cfg
}

func getCfgPath() string {

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Could not get current working directory: %v\nCould not get config path", err)
	}

	cfgPath := path.Join(cwd, "..", "..", "..", "config", "config.yaml")

	return cfgPath
}

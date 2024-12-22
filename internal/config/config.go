package config

import (
	"log"
	"os"
	"time"
   
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
// HTTPServer options
   Address     string        `yaml:"address" env-default:"localhost:8877"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"80s"`
}

func Load() *Config {
	configPath := os.Getenv("WCALC_CONFIG_PATH")
	if configPath == "" {
		log.Printf("environment variable WCALC_CONFIG_PATH  is not set, using ./config.yaml")
      configPath = "./config.yaml"
	}
	// check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return &cfg
}

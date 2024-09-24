package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ApiKey          string `yaml:"apiKey"`
	RedisUsername   string `yaml:"redisUsername"`
	RedisPassword   string `yaml:"redisPassword"`
	RedisClientAddr string `yaml:"redisClientAddr"`
	SrvListenAddr   string `yaml:"srvListenaddr"`
	MosServiceUrl   string `yaml:"mosServiceUrl"`
}

func GetConfig() *Config {
	cfgPath := "./config.yaml"
	if cfgPath == "" {
		log.Fatalf("empty cfg path %s", cfgPath)
	}

	if _, err := os.Stat(cfgPath); err != nil {
		log.Fatalf("cfg for path:%s not exists", cfgPath)
	}

	cfg := &Config{}

	if err := cleanenv.ReadConfig(cfgPath, cfg); err != nil {
		log.Fatal(err)
	}

	return cfg
}

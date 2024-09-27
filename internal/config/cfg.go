package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	ApiKey          string
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

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	cfg := &Config{
		ApiKey: os.Getenv("APIKEY"),
	}

	if err := cleanenv.ReadConfig(cfgPath, cfg); err != nil {
		log.Fatal(err)
	}

	return cfg
}

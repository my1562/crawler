package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ApiUrl string
	Redis  string
}

func NewConfig() *Config {

	envFile := ".env"

	injectedEnvFile := os.Getenv("ENV_FILE")
	if injectedEnvFile != "" {
		envFile = injectedEnvFile
	}

	err := godotenv.Load(envFile)
	if err != nil {
		log.Println(err)
	}

	config := &Config{
		ApiUrl: os.Getenv("API_URL"),
		Redis:  os.Getenv("REDIS"),
	}
	if config.Redis == "" {
		config.Redis = "127.0.0.1:6379"
	}

	return config
}

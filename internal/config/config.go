package config

import (
	"log"
	"os"
)

type Config struct {
	Env         string
	HttpAddress string
}

func MustLoad() *Config {
	env := os.Getenv("ENV")
	if env == "" {
		log.Fatal("ENV is not set")
	}

	httpAddress := os.Getenv("HTTP_ADDRESS")
	if env == "" {
		log.Fatal("HTTP_ADDRESS is not set")
	}

	return &Config{
		Env:         env,
		HttpAddress: httpAddress,
	}
}

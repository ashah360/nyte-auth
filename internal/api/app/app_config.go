package app

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	TokenTTL                   time.Duration
	Port                       string
	RequireRecaptchaToRegister bool
}

func NewConfig() *Config {
	cfg := new(Config)

	if tokenTTLSecs, ok := os.LookupEnv("TOKEN_TTL_SECS"); ok {
		tokenTTL, err := strconv.Atoi(tokenTTLSecs)
		if err != nil {
			log.Fatal(err)
		}

		cfg.TokenTTL = time.Duration(tokenTTL) * time.Second
	}

	if port, ok := os.LookupEnv("PORT"); ok {
		cfg.Port = port
	} else {
		cfg.Port = "4444"
	}

	return cfg
}

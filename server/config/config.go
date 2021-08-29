package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Env                string
	LogLevel           string
	Port               string
	DBConnectionString string
	ReadTimeout        int
	WriteTimeout       int
}

func NewConfig() *Config {
	if os.Getenv("ENV") == "" {
		os.Setenv("ENV", "development")
	}

	if os.Getenv("LOG_LEVEL") == "" {
		os.Setenv("LogLevel", "debug")
	}

	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "5000")
	}

	if os.Getenv("SERVER_READ_TIMEOUT") == "" {
		os.Setenv("SERVER_READ_TIMEOUT", "2")
	}

	if os.Getenv("SERVER_WRITE_TIMEOUT") == "" {
		os.Setenv("SERVER_WRITE_TIMEOUT", "2")
	}

	serverReadTimeout, err := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))

	if err != nil {
		log.Panicf("Invalid value %v for SERVER_READ_TIMEOUT", os.Getenv("SERVER_READ_TIMEOUT"))
	}

	serverWriteTimeout, err := strconv.Atoi(os.Getenv("SERVER_WRITE_TIMEOUT"))

	if err != nil {
		log.Panicf("Invalid value %v for SERVER_WRITE_TIMEOUT", os.Getenv("SERVER_WRITE_TIMEOUT"))
	}

	return &Config{
		Env:                os.Getenv("ENV"),
		LogLevel:           os.Getenv("LOG_LEVEL"),
		Port:               os.Getenv("PORT"),
		DBConnectionString: os.Getenv("DB_CONNECTION_STRING"),
		ReadTimeout:        serverReadTimeout,
		WriteTimeout:       serverWriteTimeout,
	}
}

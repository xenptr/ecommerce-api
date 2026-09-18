package config

import "os"

type Config struct {
	AppPort   string
	JWTSecret []byte

	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string

	RedisHost     string
	RedisPort     string
	RedisUsername string
	RedisPassword string
}

func Load() *Config {
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	return &Config{
		AppPort:   appPort,
		JWTSecret: []byte(os.Getenv("JWT_SECRET")),

		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),
		DBUser: os.Getenv("DB_USER"),
		DBPass: os.Getenv("DB_PASSWORD"),
		DBName: os.Getenv("DB_NAME"),

		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisUsername: os.Getenv("REDIS_USERNAME"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
	}
}

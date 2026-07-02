package config

import (
	"log"
	"movie_ticket/utils"
)

type Config struct {
	Port      string
	JWTSecret string

	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

func LoadConfig() *Config {
	jwtSecret := utils.GetEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set in environment variables")
	}
	return &Config{
		Port:      utils.GetEnv("PORT", "9090"),
		JWTSecret: jwtSecret,

		DBUser:     utils.GetEnv("DB_USER", "root"),
		DBPassword: utils.GetEnv("DB_PASSWORD", "root"),
		DBHost:     utils.GetEnv("DB_HOST", "127.0.0.1"),
		DBPort:     utils.GetEnv("DB_PORT", "3306"),
		DBName:     utils.GetEnv("DB_NAME", "movie_ticket_api"),
	}
}

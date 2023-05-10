package config

import (
	"os"

	"github.com/joho/godotenv"
)

type ConfigStruct struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	EmailPassword string
}

var Config ConfigStruct

func LoadConfig() {
	godotenv.Load(".env")

	Config = ConfigStruct{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		EmailPassword: os.Getenv("EMAIL_PASSWORD"),
	}
}

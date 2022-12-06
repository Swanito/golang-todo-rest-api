package server

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func loadEnvVars() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func GetApp() App {
	loadEnvVars()
	return App{
		HTTPPort: os.Getenv("HTTP_PORT"),
		DBConfig: DBConfig{
			Username: os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
		},
	}
}

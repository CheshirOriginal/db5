package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser string
	DBPass string
	DBHost string
	DBPort string
	DBName string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env не найден, продолжаем без него")
	}

	cfg := Config{
		DBUser: os.Getenv("DB_USER"),
		DBPass: os.Getenv("DB_PASS"),
		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),
		DBName: os.Getenv("DB_NAME"),
	}

	if cfg.DBUser == "" || cfg.DBPass == "" {
		log.Fatal("DB_USER или DB_PASS не заданы в .env")
	}

	return cfg
}

func (c Config) GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
}

package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Pass     string
	Name     string
	SSLMode  string
	TimeZone string
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

type CloudinaryConfig struct {
	Name      string
	ApiKey    string
	ApiSecret string
}

type Config struct {
	Env          string
	Port         string
	AppURL       string
	GroqAPIKey   string
	DuitkuAPIKey string
	LogLevel     string
	JWT          JWTConfig
	Cloudinary   CloudinaryConfig

	DB DBConfig
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env")
	}

	accessTTLHours, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRATION_HOURS"))
	refreshTTLHours, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_EXPIRATION_HOURS"))

	return &Config{
		Env:          os.Getenv("ENV"),
		Port:         os.Getenv("PORT"),
		AppURL:       os.Getenv("APP_URL"),
		GroqAPIKey:   os.Getenv("GROQ_API_KEY"),
		DuitkuAPIKey: os.Getenv("DUITKU_API_KEY"),
		LogLevel:     os.Getenv("LOG_LEVEL"),

		JWT: JWTConfig{
			AccessSecret:  os.Getenv("JWT_ACCESS_SECRET"),
			RefreshSecret: os.Getenv("JWT_REFRESH_SECRET"),
			AccessTTL:     time.Duration(accessTTLHours) * time.Hour,
			RefreshTTL:    time.Duration(refreshTTLHours) * time.Hour,
		},

		Cloudinary: CloudinaryConfig{
			Name:      os.Getenv("CLOUDINARY_NAME"),
			ApiKey:    os.Getenv("CLOUDINARY_API_KEY"),
			ApiSecret: os.Getenv("CLOUDINARY_API_SECRET"),
		},

		DB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Pass:     os.Getenv("DB_PASS"),
			Name:     os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
			TimeZone: os.Getenv("DB_TIMEZONE"),
		},
	}
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s timezone=%s",
		c.User, c.Pass, c.Host, c.Port, c.Name, c.SSLMode, c.TimeZone,
	)
}

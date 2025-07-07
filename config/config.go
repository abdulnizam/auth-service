package config

import (
	"os"
)

type Config struct {
	Port           string
	DBUser         string
	DBPass         string
	DBHost         string
	DBPort         string
	DBName         string
	JWTSecret      string
	GoogleClientID string
	GoogleSecret   string
	GoogleRedirect string
	SMTPHost       string
	SMTPPort       string
	SMTPUser       string
	SMTPPass       string
	EmailFrom      string
}

func LoadConfig() Config {
	return Config{
		Port:           getEnv("PORT", "8080"),
		DBUser:         getEnv("DB_USER", "root"),
		DBPass:         getEnv("DB_PASS", "password"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "3306"),
		DBName:         getEnv("DB_NAME", "auth_service"),
		JWTSecret:      getEnv("JWT_SECRET", "supersecretkey"),
		GoogleClientID: getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleSecret:   getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirect: getEnv("GOOGLE_REDIRECT_URI", ""),
		SMTPHost:       getEnv("SMTP_HOST", ""),
		SMTPPort:       getEnv("SMTP_PORT", "587"),
		SMTPUser:       getEnv("SMTP_USER", ""),
		SMTPPass:       getEnv("SMTP_PASS", ""),
		EmailFrom:      getEnv("EMAIL_FROM", "no-reply@example.com"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

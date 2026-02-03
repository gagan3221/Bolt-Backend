package config

import (
	"log"
	"os"
)

type Config struct {
	MongoURI     string
	DatabaseName string
	Port         string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	dbName := getEnv("DATABASE_NAME", "boltdb")
	port := getEnv("PORT", "3000")

	return &Config{
		MongoURI:     mongoURI,
		DatabaseName: dbName,
		Port:         port,
	}
}

// getEnv gets environment variable with a default fallback
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("Using default value for %s: %s", key, defaultValue)
		return defaultValue
	}
	return value
}

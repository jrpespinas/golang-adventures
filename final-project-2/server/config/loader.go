package configuration

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// load .env file
func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

// Return environment variables from .env file
func Getenv(key string) string {
	return os.Getenv(key)
}

// Return default PORT number
func GetPort(port string) string {
	if port == "" {
		return ":8000"
	} else {
		return port
	}
}

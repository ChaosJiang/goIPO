package conf

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvVariable(key string) string {
	// Load the .env file in the current directory
	err := godotenv.Load(".env")
	if err !=nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}
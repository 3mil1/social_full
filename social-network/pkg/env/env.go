package env

import (
	"github.com/joho/godotenv"
	"os"
	"social-network/pkg/logger"
)

func GoDotEnvVariable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		logger.ErrorLogger.Println("Error loading .env file. Please create .env file \nwith REFRESH_SECRET \nand ACCESS_SECRET")
		panic(err)
	}

	return os.Getenv(key)
}

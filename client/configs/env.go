package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvHTTPPort() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("HTTP_PORT")
}

func EnvGRPCServerPort() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("GRPC_SERVER_PORT")
}

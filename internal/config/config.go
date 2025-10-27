package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func loadEnv() {
	_ = godotenv.Load(".env")
}

func GetAWSRegion() (string, error) {
	loadEnv()
	v := os.Getenv("AWS_REGION")
	if v == "" {
		return "", fmt.Errorf("AWS_REGION not set")
	}
	return v, nil
}

func GetAWSEndpoint() (string, error) {
	loadEnv()
	v := os.Getenv("AWS_ENDPOINT")
	if v == "" {
		return "", fmt.Errorf("AWS_ENDPOINT not set")
	}
	return v, nil
}

func GetQueueName() (string, error) {
	loadEnv()
	v := os.Getenv("QUEUE_NAME")
	if v == "" {
		return "", fmt.Errorf("QUEUE_NAME not set")
	}
	return v, nil
}

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/joho/godotenv"
)

type Message struct {
	ID   string `json:"id"`
	Body string `json:"body"`
}

type Config struct {
	AWSRegion   string
	AWSEndpoint string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load(".env")
	cfg := &Config{
		AWSRegion:   os.Getenv("AWS_REGION"),
		AWSEndpoint: os.Getenv("AWS_ENDPOINT"),
	}
	if cfg.AWSRegion == "" || cfg.AWSEndpoint == "" {
		return nil, fmt.Errorf("missing AWS_REGION or AWS_ENDPOINT in environment")
	}
	return cfg, nil
}

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(cfg.AWSRegion),
		Endpoint: aws.String(cfg.AWSEndpoint),
	})
	if err != nil {
		panic(err)
	}
	svc := sns.New(sess)

	topicName := "test-topic"
	createOut, err := svc.CreateTopic(&sns.CreateTopicInput{
		Name: aws.String(topicName),
	})
	if err != nil {
		panic(err)
	}
	topicArn := *createOut.TopicArn

	msg := Message{ID: "1", Body: "Hello from publisher!"}
	b, _ := json.Marshal(msg)
	_, err = svc.Publish(&sns.PublishInput{
		TopicArn: aws.String(topicArn),
		Message:  aws.String(string(b)),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Published message to SNS topic.")
}

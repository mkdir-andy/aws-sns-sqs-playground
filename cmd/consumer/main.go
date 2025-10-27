package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/joho/godotenv"
)

type Message struct {
    ID   string `json:"id"`
    Body string `json:"body"`
}

type Config struct {
	AWSRegion   string
	AWSEndpoint string
	QueueName   string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load(".env")
	cfg := &Config{
		AWSRegion:   os.Getenv("AWS_REGION"),
		AWSEndpoint: os.Getenv("AWS_ENDPOINT"),
		QueueName:   os.Getenv("QUEUE_NAME"),
	}
	if cfg.AWSRegion == "" || cfg.AWSEndpoint == "" || cfg.QueueName == "" {
		return nil, fmt.Errorf("missing AWS_REGION, AWS_ENDPOINT, or QUEUE_NAME in environment")
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
    svc := sqs.New(sess)

    // Get queue URL
    qOut, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{QueueName: aws.String(cfg.QueueName)})
    if err != nil {
        panic(err)
    }
    queueURL := *qOut.QueueUrl

    fmt.Printf("Listening for messages on %s...\n", cfg.QueueName)
    for {
        out, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
            QueueUrl:            aws.String(queueURL),
            MaxNumberOfMessages: aws.Int64(1),
            WaitTimeSeconds:     aws.Int64(10),
        })
        if err != nil {
            fmt.Println("Error receiving message:", err)
            time.Sleep(2 * time.Second)
            continue
        }
        for _, m := range out.Messages {
            fmt.Printf("[Consumer %s] Raw SQS message body: %s\n", cfg.QueueName, *m.Body)
            // SNS-to-SQS envelope
            var envelope struct {
                Message string `json:"Message"`
            }
            if err := json.Unmarshal([]byte(*m.Body), &envelope); err != nil {
                fmt.Printf("[Consumer %s] Failed to parse envelope: %v\n", cfg.QueueName, err)
                continue
            }
            var msg Message
            if err := json.Unmarshal([]byte(envelope.Message), &msg); err != nil {
                fmt.Printf("[Consumer %s] Failed to parse message: %v\n", cfg.QueueName, err)
                continue
            }
            fmt.Printf("[Consumer %s] Received: %+v\n", cfg.QueueName, msg)
            // Delete message
            _, _ = svc.DeleteMessage(&sqs.DeleteMessageInput{
                QueueUrl:      aws.String(queueURL),
                ReceiptHandle: m.ReceiptHandle,
            })
        }
    }
}

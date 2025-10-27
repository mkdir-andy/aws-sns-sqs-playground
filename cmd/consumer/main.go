package main

import (
	"encoding/json"
	"fmt"
	"time"

	"aws-sns-sqs-playground/internal/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Message struct {
	ID   string `json:"id"`
	Body string `json:"body"`
}

func main() {
	region, err := config.GetAWSRegion()
	if err != nil {
		panic(err)
	}
	endpoint, err := config.GetAWSEndpoint()
	if err != nil {
		panic(err)
	}
	queueName, err := config.GetQueueName()
	if err != nil {
		panic(err)
	}

	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(endpoint),
	})
	if err != nil {
		panic(err)
	}
	svc := sqs.New(sess)

	// Get queue URL
	qOut, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{QueueName: aws.String(queueName)})
	if err != nil {
		panic(err)
	}
	queueURL := *qOut.QueueUrl

	fmt.Printf("Listening for messages on %s...\n", queueName)
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
			fmt.Printf("[Consumer %s] Raw SQS message body: %s\n", queueName, *m.Body)
			// SNS-to-SQS envelope
			var envelope struct {
				Message string `json:"Message"`
			}
			if err := json.Unmarshal([]byte(*m.Body), &envelope); err != nil {
				fmt.Printf("[Consumer %s] Failed to parse envelope: %v\n", queueName, err)
				continue
			}
			var msg Message
			if err := json.Unmarshal([]byte(envelope.Message), &msg); err != nil {
				fmt.Printf("[Consumer %s] Failed to parse message: %v\n", queueName, err)
				continue
			}
			fmt.Printf("[Consumer %s] Received: %+v\n", queueName, msg)
			// Delete message
			_, _ = svc.DeleteMessage(&sqs.DeleteMessageInput{
				QueueUrl:      aws.String(queueURL),
				ReceiptHandle: m.ReceiptHandle,
			})
		}
	}
}

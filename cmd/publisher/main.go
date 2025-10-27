package main

import (
	"encoding/json"
	"fmt"

	"aws-sns-sqs-playground/internal/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
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

	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(endpoint),
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

#!/bin/bash
set -e

# Wait for Localstack to be ready
until awslocal sns list-topics; do
  echo "Waiting for Localstack..."
  sleep 2
done

# Create SNS topic
topic_arn=$(awslocal sns create-topic --name test-topic --output text --query 'TopicArn')
echo "Created SNS topic: $topic_arn"

# Create SQS queues
queue1_url=$(awslocal sqs create-queue --queue-name test-queue-1 --output text --query 'QueueUrl')
queue2_url=$(awslocal sqs create-queue --queue-name test-queue-2 --output text --query 'QueueUrl')
echo "Created SQS queues: $queue1_url, $queue2_url"

# Get ARNs
queue1_arn=$(awslocal sqs get-queue-attributes --queue-url "$queue1_url" --attribute-name QueueArn --output text --query 'Attributes.QueueArn')
queue2_arn=$(awslocal sqs get-queue-attributes --queue-url "$queue2_url" --attribute-name QueueArn --output text --query 'Attributes.QueueArn')

# Subscribe SQS queues to SNS topic
awslocal sns subscribe --topic-arn "$topic_arn" --protocol sqs --notification-endpoint "$queue1_arn"
awslocal sns subscribe --topic-arn "$topic_arn" --protocol sqs --notification-endpoint "$queue2_arn"
echo "Subscribed both SQS queues to SNS topic."
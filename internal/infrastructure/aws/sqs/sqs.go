package sqs

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/chunnior/users/pkg/config"
)

type SQS struct {
	client *sqs.SQS
}

func NewSQS(config *config.Config) (*SQS, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile:           config.AwsProfile,
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	svc := sqs.New(sess)

	return &SQS{
		client: svc,
	}, nil
}

func (s *SQS) SendMessage(queueURL string, message interface{}) error {
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	_, err = s.client.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(messageJSON)),
		QueueUrl:    aws.String(queueURL),
	})
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

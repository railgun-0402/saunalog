package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"log"
	"net/url"
	"os"
)

func HandleRequest(ctx context.Context, event events.S3Event) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := sfn.NewFromConfig(cfg)
	stateMachineArn := os.Getenv("STATE_MACHINE_ARN")

	// 複数ファイルがアップロードされる可能性から繰り返し処理
	for _, record := range event.Records {
		log.Printf("failed to decode key: %v", record)
		s3 := record.S3

		// decode object key
		rawKey := s3.Object.Key
		key, err := url.QueryUnescape(rawKey) // URLエンコードされているとファイル名を正確に取得できない可能性
		if err != nil {
			log.Printf("failed to decode key: %v", err)
			key = rawKey
		}

		input := map[string]string{
			"fileName": key,
		}
		inputJSON, _ := json.Marshal(input)

		_, err = client.StartExecution(ctx, &sfn.StartExecutionInput{
			StateMachineArn: aws.String(stateMachineArn),
			Input:           aws.String(string(inputJSON)),
		})
		if err != nil {
			log.Printf("❌ Failed to start execution: %v", err)
			return err
		}
		log.Printf("✅ Started execution for: %s", key)
	}
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}

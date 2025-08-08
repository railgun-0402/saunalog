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
	"os"
)

func HandleRequest(ctx context.Context, event events.S3Event) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := sfn.NewFromConfig(cfg)
	stateMachineArn := os.Getenv("STATE_MACHINE_ARN")
	files := []map[string]string{}
	bucket := ""

	for _, record := range event.Records {
		log.Printf("✅ files append: %s", record.S3.Object.Key)
		// 複数ファイルがアップロードされる可能性から繰り返し処理
		bucket = record.S3.Bucket.Name
		files = append(files, map[string]string{
			"fileName": record.S3.Object.Key,
		})
	}

	// ファイルは名前も拡張子もバラバラなのでinterfaceを使用
	input := map[string]interface{}{
		"bucket": bucket,
		"files":  files,
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
	log.Printf("✅ Started execution for : %s", files)

	return nil
}

func main() {
	lambda.Start(HandleRequest)
}

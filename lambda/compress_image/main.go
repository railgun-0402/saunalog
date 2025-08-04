package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/url"
)

func HandleRequest(ctx context.Context, evnet events.S3Event) error {
	for _, record := range evnet.Records {
		s3 := record.S3
		fileName := s3.Object.Key

		// URLエンコードされてる場合があるのでデコード
		key, err := url.QueryUnescape(fileName)
		if err != nil {
			log.Printf("failed to decode S3 key: %v", err)
			key = fileName
		}

		fmt.Printf("🎯 Uploaded file: %s (bucket: %s)\n", key, s3.Bucket.Name)
	}
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}

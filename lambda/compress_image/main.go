package main

import (
	"bytes"
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"image"
	"image/jpeg"
	"log"
)

type Event struct {
	FileName string `json:"fileName"`
	Bucket   string `json:"bucket"`
}

func HandleRequest(ctx context.Context, event Event) error {
	log.Printf("Start processing: bucket=%s, key=%s", event.Bucket, event.FileName)

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("config error: %v", err)
		return err
	}
	client := s3.NewFromConfig(cfg)

	// 画像取得
	output, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(event.Bucket),
		Key:    aws.String(event.FileName),
	})
	if err != nil {
		log.Printf("get error: %v", err)
		return err
	}
	defer output.Body.Close()

	// 画像をデコード
	img, _, err := image.Decode(output.Body)
	if err != nil {
		log.Printf("decode error: %v", err)
		return err
	}

	// JPEG圧縮
	var compressed bytes.Buffer
	opts := jpeg.Options{Quality: 70}
	if err := jpeg.Encode(&compressed, img, &opts); err != nil {
		log.Printf("encode error: %v", err)
		return err
	}

	// アップロード
	newKey := "compressed/" + event.FileName
	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(event.Bucket),
		Key:    aws.String(newKey),
		Body:   bytes.NewReader(compressed.Bytes()),
	})

	if err != nil {
		log.Printf("put error: %v", err)
		return err
	}

	log.Printf("Successfully uploaded to %s/%s", event.Bucket, newKey)
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}

package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func handler(ctx context.Context, s3Event events.S3Event) {
	session := session.Must(session.NewSession())
	for _, record := range s3Event.Records {
		s3 := record.S3
		// fmt.Printf("[%s - %s] Bucket = %s, Key = %s \n", record.EventSource, record.EventTime, s3.Bucket.Name, s3.Object.Key)
		handle_file(session, s3.Bucket.Name, s3.Object.Key)
	}
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}

func handle_file(session client.ConfigProvider, bucket string, key string) {
	svc := s3.New(session)
	params := &s3.GetObjectInput{
		Bucket: aws.String(bucket), // Required
		Key:    aws.String(key),    // Required
	}
	resp, err := svc.GetObject(params)

	if err != nil {
		fmt.Printf("error: %s %s %s", bucket, key, err.Error())
		panic(err)
	}

	size := resp.ContentLength
	fmt.Printf("retrieved: %s %s %d", bucket, key, size)
}

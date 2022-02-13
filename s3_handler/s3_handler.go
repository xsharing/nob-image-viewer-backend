package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"path"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/image/draw"
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
	defer resp.Body.Close()

	if err != nil {
		fmt.Printf("error: %s %s %s\n", bucket, key, err.Error())
		panic(err)
	}

	size := resp.ContentLength
	fmt.Printf("retrieved: %s %s %d\n", bucket, key, size)

	putKeepFile(svc, bucket, key)

	srcImage, _, err := image.Decode(resp.Body)
	thumbnail := generateThumbnail(srcImage)
	encoded := new(bytes.Buffer)
	jpeg.Encode(encoded, thumbnail, &jpeg.Options{Quality: 60})

	putThumbnailImage(svc, bucket, key, encoded)

	fmt.Printf("completed: %s %s\n", bucket, key)
}

func putKeepFile(svc *s3.S3, imageBucket string, imageKey string) {
	folderStructureBucket := strings.Replace(imageBucket, "image-", "folder-structure-", 1)
	keepKey := path.Join(path.Dir(imageKey), ".keep")

	params := &s3.PutObjectInput{
		Bucket: aws.String(folderStructureBucket), // Required
		Key:    aws.String(keepKey),               // Required,
	}

	svc.PutObject(params)
}

func putThumbnailImage(svc *s3.S3, imageBucket string, imageKey string, thumbnail *bytes.Buffer) {
	thumbnailBucket := strings.Replace(imageBucket, "image-", "thumbnail-", 1)
	thumbnailKey := path.Join(hash(path.Dir(imageKey)), imageKey)

	params := &s3.PutObjectInput{
		Bucket:      aws.String(thumbnailBucket),
		Key:         aws.String(thumbnailKey),
		Body:        aws.ReadSeekCloser(bytes.NewReader(thumbnail.Bytes())),
		ContentType: aws.String("image/jpeg"),
	}

	_, err := svc.PutObject(params)

	if err != nil {
		fmt.Printf("error: %s %s %s\n", thumbnailBucket, thumbnailKey, err.Error())
		panic(err)
	}
}

func generateThumbnail(src image.Image) image.Image {
	bounds := src.Bounds()
	width := float64(bounds.Max.X - bounds.Min.X)
	height := float64(bounds.Max.Y - bounds.Min.Y)

	thumbnail_height := math.Ceil(math.Max(height*250.0/width, 250.0))
	thumbnail_width := math.Ceil(width / height * thumbnail_height)
	thumbnail := image.NewRGBA(image.Rect(0, 0, int(thumbnail_width), int(thumbnail_height)))

	draw.BiLinear.Scale(thumbnail, thumbnail.Rect, src, src.Bounds(), draw.Over, nil)

	return thumbnail
}

func hash(s string) string {
	hash := sha256.Sum256([]byte(s))
	return base58.Encode(hash[:])
}

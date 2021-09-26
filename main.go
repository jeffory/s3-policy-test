package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func GetEnvOrFallback(key string, fallback string) string {
	value := os.Getenv(key)

	if len(value) == 0 {
		return fallback
	}

	return value
}

func main() {
	port := GetEnvOrFallback("SERVER_PORT", "80")
	log.Printf("Starting server, listening on port %s", port)

	http.HandleFunc("/", HandleRequest)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	//AWSAccessKeyId := os.Getenv("AWS_ACCESS_KEY_ID")
	//AWSSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	S3BucketName := GetEnvOrFallback("S3_BUCKET_NAME", "terratest")
	//AWSRegion := GetEnvOrFallback("AWS_REGION", "ap-southeast-2")

	mainCtx := context.Background()
	cfg, err := config.LoadDefaultConfig(mainCtx)			// @TODO: Will this load up from the IAM role injected vars?

	if err != nil {
		log.Fatal(err)
	}

	testedPermissions := make(map[string]interface{})

	S3Client := s3.NewFromConfig(cfg)
	testedPermissions["ListAllBuckets"], testedPermissions["ListAllBucketsError"] = TestListAllBucketsPermission(S3Client)
	testedPermissions["PutObjectBucket"], testedPermissions["PutObjectBucketError"] = TestPutObjectPermission(S3Client, S3BucketName)

	output, _ := json.Marshal(testedPermissions)
	fmt.Fprint(w, string(output))
}

func TestListAllBucketsPermission(s3Client *s3.Client) (bool, error) {
	ctx := context.Background()
	//ctx.Deadline()		// @TODO: Set a sensible deadline in case of network issues/firewall blocking

	// List all buckets
	_, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})

	if err != nil {
		log.Println(err)
		return false, err
	}

	return true, nil
}

func TestPutObjectPermission(S3Client *s3.Client, BucketName string) (bool, error) {
	ctx := context.Background()
	//ctx.Deadline()		// @TODO: Set a sensible deadline in case of network issues/firewall blocking

	timeSuffix, _ := time.Now().MarshalText()

	f := io.LimitReader(rand.Reader, 500 * 1024)

	// @TODO: Fix the random bytes generation - either the reader needs seek functionality, or the S3 SHA256 check disabled in the SDK
	_, err := S3Client.PutObject(ctx, &s3.PutObjectInput{
		Body:   f,	  // 500 KB of Random bytes
		Bucket: aws.String(BucketName),
		Key:    aws.String(fmt.Sprintf("test-file-%s", timeSuffix)),
		ContentLength: 500 * 1024,
	})

	if err != nil {
		log.Println(err)
		return false, err
	}

	return true, nil
}

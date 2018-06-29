package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Response Json Object for responses
type S3Response struct {
	Name string `json:"name"`
	Date string `json:"created"`
}

var wg sync.WaitGroup

func createBucket(bucket string, region string) {

	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	// Create S3 service client
	svc := s3.New(sess)

	// Create the S3 Bucket
	_, err = svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		exitErrorf("Unable to create bucket %q, %v", bucket, err)
	}

	// Wait until bucket is created before finishing
	fmt.Printf("Waiting for bucket %q to be created...\n", bucket)

	err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		exitErrorf("Error occurred while waiting for bucket to be created, %v", bucket)
	}

	fmt.Printf("Bucket %q successfully created\n", bucket)
}

func listBuckets() []S3Response {

	var timeout time.Duration
	timeout = time.Duration(60) * time.Second

	// Create a context with a timeout that will abort the upload if it takes
	// more than the passed in timeout.
	ctx := context.Background()
	var cancelFn func()

	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
	}

	defer cancelFn()

	sess := session.Must(session.NewSession())
	s3svc := s3.New(sess)

	buckets := []S3Response{}
	result, err := s3svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		fmt.Println("Failed to list buckets", err)
		return buckets
	}

	for _, bucket := range result.Buckets {
		slot := S3Response{aws.StringValue(bucket.Name), bucket.CreationDate.String()}
		buckets = append(buckets, slot)
	}
	return buckets
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

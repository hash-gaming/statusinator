package main

import (
	"fmt"
	"os"

	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	_ "github.com/joho/godotenv/autoload"
)

func prettyPrint(unformattedJSON string) string {
	formattedBytes, err := json.MarshalIndent(unformattedJson, "", "  ")
	if err {
		return ""
	}

	return string(formattedBytes)
}

func main() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	sess, err := session.NewSession()
	if err != nil {
			fmt.Println("Error creating session ", err)
			return
	}
	svc := s3.New(sess, aws.NewConfig().WithRegion("us-west-2"))

	bucket, ok := os.LookupEnv("BUCKET_NAME")
	if !ok {
		fmt.Println("No bucket name found in the BUCKET_NAME env variable. Exiting.")
		os.Exit(1)
		return
	}
	input := &s3.ListObjectsV2Input{
    Bucket:  aws.String(bucket),
    MaxKeys: aws.Int64(2),
	}

	result, err := svc.ListObjectsV2(input)
	if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
					switch aerr.Code() {
					case s3.ErrCodeNoSuchBucket:
						fmt.Println(s3.ErrCodeNoSuchBucket, aerr.Error())
					default:
						fmt.Println(aerr.Error())
					}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
	}

	fmt.Println(prettyPrint(result))
}

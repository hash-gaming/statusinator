package util

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"
)

// GetAWSSession instantiates a AWS session
func GetAWSSession(options session.Options) *session.Session {
	return session.Must(session.NewSessionWithOptions(options))
}

// GetSTSClient instantiates a AWS STS client given a session and region
func GetSTSClient(sesh *session.Session, region string) *sts.STS {
	return sts.New(sesh, aws.NewConfig().WithRegion(region))
}

// GetS3Client instantiates a AWS S3 client given a session and region
func GetS3Client(sesh *session.Session, region string) *s3.S3 {
	return s3.New(sesh, aws.NewConfig().WithRegion(region))
}

// HandleS3Error handles any AWS S3 error
func HandleS3Error(err error) {
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
}

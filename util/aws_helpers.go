package util

import (
	"errors"
	"fmt"
	"os"
	"regexp"

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

// GetSessionName generates the session name based off of the username and role
func GetSessionName(roleArn string) (string, error) {
	user, _ := os.LookupEnv("USER")
	roleRegex := regexp.MustCompile("arn:aws:iam::[0-9]{12}:role/([a-zA-Z0-9-]+)")
	match := roleRegex.FindAllStringSubmatch(roleArn, -1)

	if !roleRegex.MatchString(roleArn) {
		return "", errors.New("Invalid Role ARN")
	}

	return fmt.Sprintf("%s-%s-statusinator", user, match[0][1]), nil
}

// HandleAWSError handles any AWS error
func HandleAWSError(err error, handler func(arr awserr.Error) string) string {
	if aerr, ok := err.(awserr.Error); ok {
		return handler(aerr)
	} else {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		return err.Error()
	}
}

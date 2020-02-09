package main

import (
	"fmt"
	"os"

	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"

	_ "github.com/joho/godotenv/autoload"
)

const (
	EnvBucketName     = "BUCKET_NAME"
	EnvRegion         = "REGION"
	EnvServiceRoleArn = "SERVICE_ROLE_ARN"
)

func prettyPrint(unformattedJSON interface{}) string {
	formattedBytes, err := json.MarshalIndent(unformattedJSON, "", "  ")
	if err != nil {
		return ""
	}

	return string(formattedBytes)
}

func checkEnv() {
	envVars := [3]string{
		EnvBucketName,
		EnvRegion,
		EnvServiceRoleArn,
	}

	for _, v := range envVars {
		_, ok := os.LookupEnv(v)
		if !ok {
			fmt.Println(fmt.Printf("No value found for %s in the .env file.", v))
			os.Exit(1)
			return
		}
	}
}

func getAWSSession(options session.Options) *session.Session {
	return session.Must(session.NewSessionWithOptions(options))
}

func getSTSClient(sesh *session.Session, region string) *sts.STS {
	return sts.New(sesh, aws.NewConfig().WithRegion(region))
}

func getS3Client(sesh *session.Session, region string) *s3.S3 {
	return s3.New(sesh, aws.NewConfig().WithRegion(region))
}

func main() {
	checkEnv()

	regionFromEnv, _ := os.LookupEnv(EnvRegion)
	roleArn, _ := os.LookupEnv(EnvServiceRoleArn)
	bucket, _ := os.LookupEnv(EnvBucketName)
	roleSessionName := "statusinator-test-session"

	ownAccountSesh := getAWSSession(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	stsClient := getSTSClient(ownAccountSesh, regionFromEnv)

	serviceAssumeRoleInput := &sts.AssumeRoleInput{
		RoleArn:         &roleArn,
		RoleSessionName: &roleSessionName,
	}

	assumedRole, assumeRoleErr := stsClient.AssumeRole(serviceAssumeRoleInput)
	if assumeRoleErr != nil {
		fmt.Println(assumeRoleErr)
	}

	cloudAccountSesh := getAWSSession(session.Options{
		Config: *aws.NewConfig().WithCredentials(
			credentials.NewStaticCredentials(
				*assumedRole.Credentials.AccessKeyId,
				*assumedRole.Credentials.SecretAccessKey,
				*assumedRole.Credentials.SessionToken,
			),
		),
	})

	s3Client := getS3Client(cloudAccountSesh, regionFromEnv)

	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucket),
		MaxKeys: aws.Int64(2),
	}

	result, err := s3Client.ListObjectsV2(input)
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

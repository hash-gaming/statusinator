package main

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/YashdalfTheGray/statusinator/auth"
	"github.com/YashdalfTheGray/statusinator/env"
	"github.com/YashdalfTheGray/statusinator/util"
)

func handleS3Error(err error) {
	s3Handler := func(arr awserr.Error) string {
		switch arr.Code() {
		case s3.ErrCodeNoSuchBucket:
			return s3.ErrCodeNoSuchBucket + arr.Error()
		default:
			return arr.Error()
		}
	}

	fmt.Println(util.HandleAWSError(err, s3Handler))
}

func main() {
	env.Check()

	var creds *credentials.Credentials
	roleArn := env.Get(env.ServiceRoleArn)

	if env.IsDevelopment() {
		assumedRole, assumeRoleErr := auth.GetDevAuth(roleArn)
		if assumeRoleErr != nil {
			fmt.Println(assumeRoleErr)
		}

		creds = credentials.NewStaticCredentials(
			*assumedRole.Credentials.AccessKeyId,
			*assumedRole.Credentials.SecretAccessKey,
			*assumedRole.Credentials.SessionToken,
		)
	} else {
		fmt.Println("No other runtime environment currently configured besides development.")
		os.Exit(1)
		return
	}

	cloudAccountSesh := util.GetAWSSession(session.Options{
		Config: *aws.NewConfig().WithCredentials(creds),
	})

	s3Client := util.GetS3Client(cloudAccountSesh, env.Get(env.Region))

	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(env.Get(env.BucketName)),
		MaxKeys: aws.Int64(2),
	}

	result, err := s3Client.ListObjectsV2(input)
	if err != nil {
		handleS3Error(err)
	}

	fmt.Println(util.PrettyPrint(result))
}

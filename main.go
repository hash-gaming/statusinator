package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"

	_ "github.com/joho/godotenv/autoload"

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

	region, _ := os.LookupEnv(env.Region)
	roleArn, _ := os.LookupEnv(env.ServiceRoleArn)
	bucketName, _ := os.LookupEnv(env.BucketName)
	roleSessionName := "statusinator-test-session"

	ownAccountSesh := util.GetAWSSession(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	stsClient := util.GetSTSClient(ownAccountSesh, region)

	serviceAssumeRoleInput := &sts.AssumeRoleInput{
		RoleArn:         &roleArn,
		RoleSessionName: &roleSessionName,
	}

	assumedRole, assumeRoleErr := stsClient.AssumeRole(serviceAssumeRoleInput)
	if assumeRoleErr != nil {
		fmt.Println(assumeRoleErr)
	}

	cloudAccountSesh := util.GetAWSSession(session.Options{
		Config: *aws.NewConfig().WithCredentials(
			credentials.NewStaticCredentials(
				*assumedRole.Credentials.AccessKeyId,
				*assumedRole.Credentials.SecretAccessKey,
				*assumedRole.Credentials.SessionToken,
			),
		),
	})

	s3Client := util.GetS3Client(cloudAccountSesh, region)

	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucketName),
		MaxKeys: aws.Int64(2),
	}

	result, err := s3Client.ListObjectsV2(input)
	if err != nil {
		handleS3Error(err)
	}

	fmt.Println(util.PrettyPrint(result))
}

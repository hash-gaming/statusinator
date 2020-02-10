package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sts"

	_ "github.com/joho/godotenv/autoload"

	"github.com/YashdalfTheGray/statusinator/util"
)

func main() {
	util.CheckEnv()

	regionFromEnv, _ := os.LookupEnv(util.EnvRegion)
	roleArn, _ := os.LookupEnv(util.EnvServiceRoleArn)
	bucketName, _ := os.LookupEnv(util.EnvBucketName)
	roleSessionName := "statusinator-test-session"

	ownAccountSesh := util.GetAWSSession(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	stsClient := util.GetSTSClient(ownAccountSesh, regionFromEnv)

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

	s3Client := util.GetS3Client(cloudAccountSesh, regionFromEnv)

	input := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucketName),
		MaxKeys: aws.Int64(2),
	}

	result, err := s3Client.ListObjectsV2(input)
	util.HandleS3Error(err)

	fmt.Println(util.PrettyPrint(result))
}

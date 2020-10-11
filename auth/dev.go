package auth

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/hash-gaming/statusinator/env"
	"github.com/hash-gaming/statusinator/util"
)

// GetDevAuth returns the credentials for an authenticated session using AWS STS
func GetDevAuth(roleArn string) (*sts.AssumeRoleOutput, error) {
	roleSessionName, err := util.GetSessionName(roleArn)
	if err != nil {
		return nil, err
	}

	ownAccountSesh := util.GetAWSSession(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})

	stsClient := util.GetSTSClient(ownAccountSesh, env.Get(env.Region))

	serviceAssumeRoleInput := &sts.AssumeRoleInput{
		RoleArn:         &roleArn,
		RoleSessionName: &roleSessionName,
	}

	return stsClient.AssumeRole(serviceAssumeRoleInput)
}

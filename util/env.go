package util

import (
	"fmt"
	"os"
)

// Environment variable constants so that we don't keep messing stuff up
const (
	EnvBucketName     = "BUCKET_NAME"
	EnvRegion         = "REGION"
	EnvServiceRoleArn = "SERVICE_ROLE_ARN"
)

// CheckEnv ensures that the necessary environment variables are present
func CheckEnv() {
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

package env

import (
	"fmt"
	"os"
)

// Environment variable constants so that we don't keep messing stuff up
const (
	BucketName     = "BUCKET_NAME"
	Region         = "REGION"
	ServiceRoleArn = "SERVICE_ROLE_ARN"
)

// Check ensures that the necessary environment variables are present
func Check() {
	envVars := [3]string{
		BucketName,
		Region,
		ServiceRoleArn,
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

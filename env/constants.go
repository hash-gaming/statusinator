package env

// Environment variable constants so that we don't keep messing stuff up
const (
	BucketName     = "BUCKET_NAME"
	Region         = "REGION"
	ServiceRoleArn = "SERVICE_ROLE_ARN"
)

var allEnvKeys = []string{
	BucketName,
	Region,
	ServiceRoleArn,
}

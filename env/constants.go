package env

// Environment variable constants so that we don't keep messing stuff up
const (
	// Required Environment Variable Keys
	BucketName     = "BUCKET_NAME"
	Region         = "REGION"
	RuntimeEnv     = "ENV"
	ServiceRoleArn = "SERVICE_ROLE_ARN"

	// Optional Environment Variable Keys
	AWSAccessKeyID     = "AWS_ACCESS_KEY_ID"
	AWSSecretAccessKey = "AWS_SECRET_ACCESS_KEY"

	// Valid values for the ENV environment variable
	RuntimeDev  = "development"
	RuntimeProd = "production"
)

var requiredEnvKeys = []string{
	BucketName,
	Region,
	RuntimeEnv,
	ServiceRoleArn,
}

var optionalEnvKeys = []string{
	AWSAccessKeyID,
	AWSSecretAccessKey,
}

var envVarMap = make(map[string]string, len(requiredEnvKeys)+len(optionalEnvKeys))

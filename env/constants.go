package env

// Environment variable constants so that we don't keep messing stuff up
const (
	// Environment variable keys
	BucketName     = "BUCKET_NAME"
	Region         = "REGION"
	RuntimeEnv     = "ENV"
	ServiceRoleArn = "SERVICE_ROLE_ARN"

	// Valid values for the ENV environment variable
	RuntimeDev  = "development"
	RuntimeProd = "production"
)

var allEnvKeys = []string{
	BucketName,
	Region,
	RuntimeEnv,
	ServiceRoleArn,
}

var envVarMap = make(map[string]string, len(allEnvKeys))

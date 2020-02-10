package env

import "os"

var envVarMap = make(map[string]string, len(allEnvKeys))

// Get returns the value of an environment variable. If it doesn't
// exist, it asks os.LookupEnv and caches the value.
func Get(envKey string) string {
	if envVarMap[envKey] == "" {
		envValue, _ := os.LookupEnv(envKey)
		envVarMap[envKey] = envValue
	}
	return envVarMap[envKey]
}

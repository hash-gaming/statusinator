package env

import (
	"fmt"
	"os"
	"strings"
)

// Check ensures that the necessary environment variables are present
func Check() {
	for _, key := range requiredEnvKeys {
		value, ok := os.LookupEnv(key)
		if !ok {
			fmt.Println(fmt.Printf("No value found for %s in the .env file.", key))
			os.Exit(1)
		}

		envVarMap[key] = value

		if key == RuntimeEnv {
			if strings.ToLower(value) != RuntimeDev && strings.ToLower(value) != RuntimeProd {
				fmt.Println("ENV can only be set to either 'development' or 'production'.")
				os.Exit(1)
			}
		}
	}
}

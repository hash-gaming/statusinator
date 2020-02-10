package env

import (
	"fmt"
	"os"
)

// Check ensures that the necessary environment variables are present
func Check() {
	for _, v := range allEnvKeys {
		_, ok := os.LookupEnv(v)
		if !ok {
			fmt.Println(fmt.Printf("No value found for %s in the .env file.", v))
			os.Exit(1)
			return
		}
	}
}

package util

import "encoding/json"

// PrettyPrint prints out something that looks like JSON
// prettily
func PrettyPrint(unformattedJSON interface{}) string {
	formattedBytes, err := json.MarshalIndent(unformattedJSON, "", "  ")
	if err != nil {
		return ""
	}

	return string(formattedBytes)
}

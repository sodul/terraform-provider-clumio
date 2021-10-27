// Copyright 2021. Clumio, Inc.

// Contains the util functions used by the providers and resources

package clumio

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// getStringValue returns the string value of the key if present.
func getStringValue(d *schema.ResourceData, key string) string {
	value := ""
	if d.Get(key) != nil {
		value = fmt.Sprintf("%v", d.Get(key))
	}
	return value
}

//Utility function to return a string value from a map if the key exists
func getStringValueFromMap(keyVals map[string]interface{}, key string) *string{
	if v, ok := keyVals[key].(string); ok && v != "" {
		return &v
	}
	return nil
}

// Utility function to determine if it is unit or acceptance test
func isAcceptanceTest() bool{
	return os.Getenv("TF_ACC") == "true" || os.Getenv("TF_ACC") == "True" ||
		os.Getenv("TF_ACC") == "1"
}

// RequireOneOf verifies that at least one environment variable is non-empty or returns an error.
//
// If at lease one environment variable is non-empty, returns the first name and value.
func RequireOneOf(names []string, usageMessage string) (string, string, error) {
	for _, variable := range names {
		value := os.Getenv(variable)

		if value != "" {
			return variable, value, nil
		}
	}

	return "", "", fmt.Errorf("at least one environment variable of %v must be set. Usage: %s", names, usageMessage)
}
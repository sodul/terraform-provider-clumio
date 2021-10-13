// Copyright 2021. Clumio, Inc.

// Contains the util functions used by the providers and resources

package clumio

import (
	"fmt"

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
package clumio

import (
	"os"
	"testing"

	"github.com/clumio-code/terraform-provider-clumio/clumio/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ProviderFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var ProviderFactories = map[string]func() (*schema.Provider, error){
	"clumio": func() (*schema.Provider, error) {
		return New(!common.IsAcceptanceTest())(), nil
	},
}

// UtilTestAccPreCheckAws verifies and sets required provider testing configuration
//
// This PreCheck function should be present in every acceptance test. It allows
// test configurations to omit a provider configuration with region and ensures
// testing functions that attempt to call AWS APIs are previously configured.
//
// These verifications and configuration are preferred at this level to prevent
// provider developers from experiencing less clear errors for every test.
func UtilTestAccPreCheckAws(t *testing.T) {
	UtilTestFailIfAllEmpty(t,
		[]string{common.AwsAccessKeyId, awsProfile, awsSharedCredsFile},
		"One of AWS_ACCESS_KEY_ID, AWS_PROFILE or AWS_SHARED_CREDENTIALS_FILE must be set")
	if os.Getenv(common.AwsAccessKeyId) != "" {
		UtilTestFailIfEmpty(t, common.AwsSecretAccessKey, common.AwsSecretAccessKey+" cannot be empty.")
	}
}

// UtilTestAccPreCheckClumio validates that the required environment variables are set before
// the acceptance test is executed.
func UtilTestAccPreCheckClumio(t *testing.T) {
	UtilTestFailIfEmpty(t, common.ClumioTestAwsAccountId, common.ClumioTestAwsAccountId+" cannot be empty")
	UtilTestFailIfEmpty(t, common.ClumioApiToken, common.ClumioApiToken+" cannot be empty.")
	UtilTestFailIfEmpty(t, common.ClumioApiBaseUrl, common.ClumioApiBaseUrl+" cannot be empty.")
	UtilTestFailIfEmpty(t, common.AwsRegion, common.AwsRegion+" cannot be empty")
}

// UtilTestFailIfAllEmpty verifies that at least one environment variable is non-empty or fails the test.
// If at lease one environment variable is non-empty, returns the first name and value.
func UtilTestFailIfAllEmpty(t *testing.T, names []string, usageMessage string) (string, string) {
	t.Helper()

	name, value, err := common.RequireOneOf(names, usageMessage)
	if err != nil {
		t.Fatal(err)
		return "", ""
	}

	return name, value
}

// UtilTestFailIfEmpty verifies that an environment variable is non-empty or fails the test.
//
// For acceptance tests, this function must be used outside PreCheck functions to set values for configurations.
func UtilTestFailIfEmpty(t *testing.T, name string, usageMessage string) string {
	t.Helper()

	value := os.Getenv(name)

	if value == "" {
		t.Fatalf("environment variable %s must be set. Usage: %s", name, usageMessage)
	}

	return value
}

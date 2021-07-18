// Copyright 2021. Clumio, Inc.

// Acceptance tests for provider.
package clumio

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"clumio": func() (*schema.Provider, error) {
		return New(true)(), nil
	},
}

func TestProvider(t *testing.T) {
	if err := New(true)().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}


// testAccPreCheck verifies and sets required provider testing configuration
//
// This PreCheck function should be present in every acceptance test. It allows
// test configurations to omit a provider configuration with region and ensures
// testing functions that attempt to call AWS APIs are previously configured.
//
// These verifications and configuration are preferred at this level to prevent
// provider developers from experiencing less clear errors for every test.
func testAccPreCheck(t *testing.T) {

}

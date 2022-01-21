// Copyright 2021. Clumio, Inc.

// Acceptance tests for provider.
package clumio

import (
	"testing"
)

func TestProvider(t *testing.T) {
	if err := New(true)().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

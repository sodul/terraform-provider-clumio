// Copyright 2021. Clumio, Inc.

// Acceptance test for clumio_organizational_unit resource.
package clumio_organizational_unit_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/clumio-code/terraform-provider-clumio/clumio"
	"github.com/clumio-code/terraform-provider-clumio/clumio/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceClumioOrganizationalUnit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { clumio.UtilTestAccPreCheckClumio(t) },
		ProviderFactories: clumio.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: getTestAccResourceClumioOrganizationalUnit(false),
			},
			{
				Config: getTestAccResourceClumioOrganizationalUnit(true),
			},
		},
	})
}

func getTestAccResourceClumioOrganizationalUnit(update bool) string {
	baseUrl := os.Getenv(common.ClumioApiBaseUrl)
	name := "acceptance-test-ou"
	if update {
		name = "acceptance-test-ou-updated"
	}
	return fmt.Sprintf(testAccResourceClumioOrganizationalUnit, baseUrl, name)
}

const testAccResourceClumioOrganizationalUnit = `
provider clumio{
   clumio_api_base_url = "%s"
}

resource "clumio_organizational_unit" "test_ou" {
  name = "%s"
}
`

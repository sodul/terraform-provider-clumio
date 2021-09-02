// Copyright 2021. Clumio, Inc.

// Acceptance test for resource_clumio_callback.
package clumio

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceClumioCallback(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceClumioCallback,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"clumio_callback_resource.test", "sns_topic",
						regexp.MustCompile("sns_topic")),
				),
			},
		},
	})
}

const testAccResourceClumioCallback = `
resource "clumio_callback_resource" "test" {
  sns_topic = "clumio_test_sns_topic"
  token = "89130d52-67cb-4752-a896-730cf067aeb1"
  role_external_id = "Aasdfjhg8943kbnlasdklhkljghlkash7892r"
  account_id = "482567874266"
  region = "us-west-2"
  role_id = "TestRole-us-west-2-89130d52-67cb-4752-a896-730cf067aeb1"
  role_arn = "arn:aws:iam::482567874266:role/clumio/TestRole-us-west-2-89130d52-67cb-4752-a896-730cf067aeb1"
  clumio_event_pub_id = "arn:aws:sns:us-west-2:482567874266:ClumioInventoryTopic_89130d52-67cb-4752-a896-730cf067aeb1"
  type = "service"
  bucket_name = "clumio_terraform_test_bucket"
  canonical_user = "canonical_user"
  config_version = "1"
  discover_version = "3"
  protect_config_version = "18"
  protect_ebs_version = "19"
  protect_rds_version = "18"
  protect_ec2_mssql_version = "1"
  protect_warm_tier_version = "2"
  protect_warm_tier_dynamodb_version = "2"
}
`

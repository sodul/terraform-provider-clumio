resource "clumio_callback_resource" "example" {
  # example configuration here
  topic               = "mytopic"
  token               = "mytoken"
  role_external_id    = "role_external_id"
  account_id          = "account_id"
  region              = "region"
  role_id             = "role_id"
  role_arn            = "role_arn"
  clumio_event_pub_id = "clumio_event_pub_id"
  type                = "type"
  properties = {
    "prop1" : {
      "key1" : val1
    }
  }
  config_version                     = "1"
  discover_version                   = "3"
  protect_config_version             = "18"
  protect_ebs_version                = "19"
  protect_rds_version                = "18"
  protect_ec2_mssql_version          = "1"
  protect_warm_tier_version          = "2"
  protect_warm_tier_dynamodb_version = "2"
  protect_s3_version                 = "1"
}

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
}

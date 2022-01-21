resource "clumio_policy_rule" "example" {
  name = "example-policy-rule"
  policy_id = "policy_id"
  before_rule_id = "rule_id"
  condition = "{\"aws_account_native_id\":{\"$in\":[\"aws_account_id\", \"aws_account_id_2\"]}, \"aws_tag\":{\"$eq\":{\"key\":\"aws_tag_key\", \"value\":\"aws_tag_value\"}}}"
}

resource "aws_dynamodb_table" "hal" {
  "attribute" {
    name = "groupId"
    type = "S"
  }
  hash_key = "groupId"
  name = "hal"
  read_capacity = 0
  write_capacity = 0
  billing_mode = "PAY_PER_REQUEST"
  point_in_time_recovery {
    enabled = true
  }

  lifecycle {
    prevent_destroy = true
  }
}
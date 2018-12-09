resource "aws_ssm_parameter" "domain" {
  name = "domain"
  type = "String"
  value = "${var.domain}"
}

resource "aws_ssm_parameter" "telegram" {
  name = "telegram-key"
  type = "String"
  value = "${var.telegram-key}"
}

resource "aws_ssm_parameter" "error-group" {
  name = "error-group"
  type = "String"
  value = "${var.error-group}"
}
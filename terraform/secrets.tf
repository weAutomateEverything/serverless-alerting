provider "aws" {
  profile = "card"
  region = "eu-west-1"
}

terraform {
  backend "s3" {
    bucket = "cardpayments-states"
    key = "serverless-alerting"
    region = "eu-west-1"
    profile = "card"
  }
}

variable "domain" {
  default = "api.carddevops.co.za"
}

variable "telegram-key" {
  default = "687161264:AAGoDYyJnVgpN_7Qe126ILKfqY0I7gHWroY"
}
variable "error-group" {
  default = "1083240170"
}

variable "github_key" {
  default = "2c5f91ec0743c6bfcd6470a7d8f4e6c60a2db8d4"
}

variable "region" {
  default = "eu-central-1"
}

variable "path_to_archive" {
    default = "~/zip/ec2-recovery-lambda.zip"
}
variable "name_of_app" {
    default = "ec2-recovery-lambda"
}

variable "slack_enable" {
    default = "false"
}

variable "slack_channel" {
    default = ""
}
variable "slack_auth_token" {
    default = ""
}
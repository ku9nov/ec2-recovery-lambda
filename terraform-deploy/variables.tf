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

variable "notify_arns" {
   type        = list(string)
   description = "A list of ARNs (i.e. SNS Topic ARN) to execute when this alarm transitions into ANY state from any other state. May be overridden by the value of a more specific {alarm,ok,insufficient_data}_actions variable. "
   default     = [""]
}
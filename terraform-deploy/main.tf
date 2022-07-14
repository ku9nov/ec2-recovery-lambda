provider "aws" {
  region     = var.region
}
locals {
  alarm_actions             = var.notify_arns
  ok_actions                = var.notify_arns
  insufficient_data_actions = var.notify_arns
}
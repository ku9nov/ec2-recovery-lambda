resource "aws_lambda_function" "go_lambda" {

  filename      = var.path_to_archive
  function_name = var.name_of_app
  role          = aws_iam_role.iam_for_go_lambda.arn
  handler       = var.name_of_app
  timeout       = 900

  source_code_hash = filebase64sha256(var.path_to_archive)

  runtime = "go1.x"

  environment {
    variables = {
      SLACK_ENABLE = var.slack_enable
      SLACK_CHANNEL_ID = var.slack_channel
      SLACK_AUTH_TOKEN = var.slack_auth_token
    }
  }
}
resource "aws_cloudwatch_event_rule" "every_five_minutes" {
    name = "every-five-minutes-${var.name_of_app}-${var.region}"
    description = "Fires every five minutes"
    schedule_expression = "rate(5 minutes)"
}

resource "aws_cloudwatch_event_target" "check_go_lambda_every_five_minutes" {
    rule = aws_cloudwatch_event_rule.every_five_minutes.name
    target_id = "${var.name_of_app}-${var.region}"
    arn = aws_lambda_function.go_lambda.arn
}

resource "aws_lambda_permission" "allow_cloudwatch_to_call_go_lambda" {
    statement_id  = "AllowExecutionFromCloudWatch"
    action        = "lambda:InvokeFunction"
    function_name = aws_lambda_function.go_lambda.function_name
    principal     = "events.amazonaws.com"
    source_arn    = aws_cloudwatch_event_rule.every_five_minutes.arn
}

resource "aws_cloudwatch_metric_alarm" "lambda_alarm" {
  alarm_name          = "${aws_lambda_function.go_lambda.function_name}_lambda_alarm"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = 1
  metric_name         = "Errors"
  namespace           = "AWS/Lambda"
  period              = 60
  statistic           = "Sum"
  threshold           = 0
  alarm_actions       = try(local.alarm_actions, false) 
  ok_actions          = try(local.ok_actions, false) 
  dimensions = {
    FunctionName = "${aws_lambda_function.go_lambda.function_name}"
  }
}
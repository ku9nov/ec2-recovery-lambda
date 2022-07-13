resource "aws_iam_role" "iam_for_go_lambda" {
  name = "go_${var.name_of_app}_${var.region}"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}


resource "aws_iam_policy" "policy_for_go_lambda" {
  name = "tf-iam-role-policy-${var.name_of_app}-${var.region}"

  policy = <<POLICY
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ec2:DescribeInstances",
                "ec2:StartInstances",
                "ec2:StopInstances",
                "ec2:DescribeInstanceStatus"
            ],
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogStream",
                "logs:CreateLogGroup",
                "logs:PutLogEvents"
            ],
            "Resource": "arn:aws:logs:*:*:*"
        }
    ]
}
POLICY
}


resource "aws_iam_role_policy_attachment" "ecs_attach" {
  role       = aws_iam_role.iam_for_go_lambda.name
  policy_arn = aws_iam_policy.policy_for_go_lambda.arn
}
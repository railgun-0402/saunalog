provider "aws" {
  region = "ap-northeast-1"
}

resource "aws_s3_bucket" "lambda_deploy" {
  bucket = "saunalog-lambda-deploy-bucket"
  force_destroy = true
}

// IAM Role
resource "aws_iam_role" "lambda_exec_role" {
  name = "saunalog-lambda-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })
}

// Create Lambda Function
resource "aws_iam_role_policy_attachment" "lambda_basic_execution" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  role = aws_iam_role.lambda_exec_role.name
}

resource "aws_lambda_function" "compress_image2" {
  function_name = "compress_image2"
  role = aws_iam_role.lambda_exec_role.arn
  handler = "bootstrap"
  runtime = "provided.al2023"

  filename = "build/compress_image.zip"
  source_code_hash = filebase64sha256("build/compress_image.zip")

  environment {
    variables = {
      STAGE = "dev"
    }
  }
}

// Step Functions
resource "aws_sfn_state_machine" "media_pipeline" {
  name = "saunalog-media-pipeline"
  role_arn = aws_iam_role.stepfn_role.arn

  definition =templatefile("${path.module}/asl/compress_state_machine.json", {
    compress_image_lambda_arn = aws_lambda_function.compress_image2.arn
  })

  logging_configuration {
    level = "ALL"
    include_execution_data = true
    log_destination        = "${aws_cloudwatch_log_group.stepfn_logs.arn}:*"
  }

  depends_on = [
    aws_cloudwatch_log_group.stepfn_logs,
    aws_iam_role_policy.stepfn_logs_policy
  ]
}

// CloudWatch
resource "aws_cloudwatch_log_group" "stepfn_logs" {
  name              = "/aws/vendedlogs/states/saunalog-media-pipeline"
  retention_in_days = 14
}

resource "aws_iam_role_policy" "stepfn_logs_policy" {
  role = aws_iam_role.stepfn_role.id
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect: "Allow",
        Action: [
          "logs:CreateLogDelivery",
          "logs:GetLogDelivery",
          "logs:UpdateLogDelivery",
          "logs:DeleteLogDelivery",
          "logs:ListLogDeliveries",
          "logs:PutResourcePolicy",
          "logs:DescribeResourcePolicies",
          "logs:DescribeLogGroups"
        ],
        Resource: "*"
      }
    ]
  })
}

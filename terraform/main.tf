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

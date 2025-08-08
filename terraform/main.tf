provider "aws" {
  region = "ap-northeast-1"
}

resource "aws_s3_bucket" "lambda_deploy" {
  bucket = "saunalog-lambda-deploy-bucket"
  force_destroy = true
}

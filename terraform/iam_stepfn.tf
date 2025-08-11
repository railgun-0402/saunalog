
// StepFunctionsのIAM Role
resource "aws_iam_role" "stepfn_role" {
  name = "saunalog-stepfn-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "states.amazonaws.com"
      }
    }]
  })
}

resource "aws_iam_role_policy" "stepfn_policy" {
  name = "saunalog-stepfn-policy"
  role = aws_iam_role.stepfn_role.name
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = ["lambda:InvokeFunction"]
      Effect = "Allow"
      Resource = [
        aws_lambda_function.compress_image2.arn
      ]
    },
      // TODO: compress_video用
//      {
//        Action = ["lambda:InvokeFunction"]
//        Effect = "Allow"
//        Resource = [
//          aws_lambda_function.compress_image2.arn
//        ]
//      }
    ]
  })
}

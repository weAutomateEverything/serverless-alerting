variable "name" {
  default = "serverless-alerting"
}

resource "aws_iam_role" "pipeline" {
  name = "pipeline-${var.name}"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "codepipeline.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}


resource "aws_s3_bucket" "build" {
  bucket = "build-${var.name}-${var.domain}"
}

resource "aws_iam_role" "codebuild" {
  name = "codebuild-${var.name}"
  assume_role_policy = <<EOF1
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "codebuild.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF1
}

resource "aws_iam_role_policy" "pipeline" {
  name = "pipeline-${var.name}"
  role = "${aws_iam_role.pipeline.name}"
  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Resource": [
        "*"
      ],
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ]
    },
    {
      "Effect": "Allow",
      "Action": [
        "s3:*"
      ],
      "Resource": [
        "${aws_s3_bucket.build.arn}",
        "${aws_s3_bucket.build.arn}/*",
        "${aws_s3_bucket.lambda.arn}",
        "${aws_s3_bucket.lambda.arn}/*"

      ]
    },
    {
        "Effect": "Allow",
        "Action": [
           "codebuild:*"
        ],
        "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Resource": [
        "*"
      ],
      "Action": "codecommit:*"
    }
  ]
}
POLICY
}

resource "aws_iam_role_policy" "build" {
  role = "${aws_iam_role.codebuild.name}"
  policy = <<POLICY
{
	"Statement": [
      {
          "Action": [
              "route53:ListHostedZones"
          ],
          "Resource": "*",
          "Effect": "Allow"
      },
	  {
			"Action": [
				"s3:*"
			],
			"Resource": "*",
			"Effect": "Allow"
		},
		{
			"Action": [
				"lambda:*"
			],
			"Resource": [
				"*"
			],
			"Effect": "Allow"
		},
		{
			"Action": [
				"apigateway:*"
			],
			"Resource": [
				"*"
			],
			"Effect": "Allow"
		},
		{
			"Action": [
				"iam:GetRole",
				"iam:CreateRole",
				"iam:DeleteRole",
				"iam:PutRolePolicy"
			],
			"Resource": [
				"*"
			],
			"Effect": "Allow"
		},
		{
			"Action": [
				"iam:AttachRolePolicy",
				"iam:DeleteRolePolicy",
				"iam:DetachRolePolicy"
			],
			"Resource": [
				"*"
			],
			"Effect": "Allow"
		},
		{
			"Action": [
				"iam:PassRole"
			],
			"Resource": [
				"*"
			],
			"Effect": "Allow"
		},
		{
			"Action": [
				"cloudformation:*"
			],
			"Resource": [
				"*"
			],
			"Effect": "Allow"
		},
		{
			"Action": [
				"codedeploy:CreateApplication",
				"codedeploy:DeleteApplication",
				"codedeploy:RegisterApplicationRevision"
			],
			"Resource": [
				"*"
			],
			"Effect": "Allow"
		},
		{
			"Action": [
				"codedeploy:CreateDeploymentGroup",
				"codedeploy:CreateDeployment",
				"codedeploy:GetDeployment"
			],
			"Resource": [
				"*"
			],
			"Effect": "Allow"
		},
		{
			"Action": [
				"codedeploy:GetDeploymentConfig"
			],
			"Resource": [
				"*"
			],
			"Effect": "Allow"
		},
        {
          "Effect": "Allow",
          "Resource": [
            "*"
          ],
          "Action": [
            "logs:*"
          ]
        },
        {
          "Action": [
            "ecr:BatchCheckLayerAvailability",
            "ecr:CompleteLayerUpload",
            "ecr:GetAuthorizationToken",
            "ecr:InitiateLayerUpload",
            "ecr:PutImage",
            "ecr:UploadLayerPart"
          ],
          "Resource": "*",
          "Effect": "Allow"
      }
	],
	"Version": "2012-10-17"
}
POLICY
}


resource "aws_codebuild_project" "serverless" {
  name = "serverless-${var.name}"
  service_role = "${aws_iam_role.codebuild.arn}"
  "artifacts" {
    type = "CODEPIPELINE"
  }
  "environment" {
    compute_type = "BUILD_GENERAL1_SMALL"
    image = "zamedic/gosls"
    type = "LINUX_CONTAINER"
    environment_variable {
      name = "Domain"
      value = "${var.domain}"
    }
  }
  "source" {
    type = "CODEPIPELINE"
  }

}

resource "aws_s3_bucket" "lambda" {
  bucket = "lambda-${var.name}-${var.domain}"
  acl = "private"
}


resource "aws_codepipeline" "lambda" {
  name = "${var.name}"
  "artifact_store" {
    location = "${aws_s3_bucket.lambda.bucket}"
    type = "S3"
  }
  role_arn = "${aws_iam_role.pipeline.arn}"
  "stage" {
    "action" {
      category = "Source"
      name = "Source"
      owner = "ThirdParty"
      provider = "GitHub"
      version = "1"
      configuration {
        Owner = "weAutomateEverything"
        Repo = "serverless-alerting"
        Branch = "master"
        OAuthToken = "${var.github_key}"
      }
      output_artifacts = ["source"]
    }
    name = "SourceCode"
  }
  stage {
    "action" {
      category = "Build"
      name = "lambda-auth"
      owner = "AWS"
      provider = "CodeBuild"
      version = "1"
      configuration {
        ProjectName = "${aws_codebuild_project.serverless.name}"
      }
      input_artifacts = [
        "source"]
      output_artifacts = [
        "lambda"]
    }
    name = "labda-build"
  },
}

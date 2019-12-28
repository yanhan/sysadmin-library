provider "aws" {
  region = var.aws_region
  version = "~> 2.0"
}

data "aws_caller_identity"  "current" {}

resource "aws_s3_bucket"  "ec2_image_builder" {
  bucket = var.s3_bucket_name
  region = var.aws_region
  tags = {
    ManagedBy = "Terraform"
  }
}

resource "aws_s3_bucket_public_access_block"  "ec2_image_builder" {
  bucket = aws_s3_bucket.ec2_image_builder.id
  block_public_acls = "true"
  block_public_policy = "true"
  ignore_public_acls = "true"
  restrict_public_buckets = "true"
}

resource "aws_licensemanager_license_configuration"  "my_own_instances" {
  name = "my-own-instances"
  description = "For my own instances"
  license_counting_type = "Instance"
  tags = {
    ManagedBy = "Terraform"
  }
}

resource "aws_sns_topic"  "ec2_image_builder" {
  name = "ec2-image-builder"
  tags = {
    ManagedBy = "Terraform"
  }
}

resource "aws_sns_topic_policy"  "ec2_image_builder" {
  arn = aws_sns_topic.ec2_image_builder.arn
  policy = <<EOF
{
  "Version": "2008-10-17",
  "Id": "__default_policy_ID",
  "Statement": [
    {
      "Sid": "__default_statement_ID",
      "Effect": "Allow",
      "Principal": {
        "AWS": "*"
      },
      "Action": [
        "SNS:GetTopicAttributes",
        "SNS:SetTopicAttributes",
        "SNS:AddPermission",
        "SNS:RemovePermission",
        "SNS:DeleteTopic",
        "SNS:Subscribe",
        "SNS:ListSubscriptionsByTopic",
        "SNS:Publish",
        "SNS:Receive"
      ],
      "Resource": "${aws_sns_topic.ec2_image_builder.arn}",
      "Condition": {
        "StringEquals": {
          "AWS:SourceOwner": "${data.aws_caller_identity.current.account_id}"
        }
      }
    }
  ]
}
EOF
}

resource "aws_iam_role"  "sample_linux_image_builder_node" {
  name = "sample-linux-image-builder-node"
  assume_role_policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "Service": "ec2.amazonaws.com"
            },
            "Action": "sts:AssumeRole"
        }
    ]
}
EOF
  description = "For EC2 nodes used for building image using EC2 Image Builder"
  tags = {
    ManagedBy = "Terraform"
  }
}

resource "aws_iam_role_policy_attachment"  "sample_linux_image_builder_node_AWSSSMManagedInstanceCore" {
  role = aws_iam_role.sample_linux_image_builder_node.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

resource "aws_iam_role_policy_attachment"  "sample_linux_image_builder_node_EC2InstanceProfileForImageBuilder" {
  role = aws_iam_role.sample_linux_image_builder_node.name
  policy_arn = "arn:aws:iam::aws:policy/EC2InstanceProfileForImageBuilder"
}

resource "aws_iam_role_policy"  "sample_linux_image_builder_node" {
  role = aws_iam_role.sample_linux_image_builder_node.name
  name = "node-policy"
  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "s3:GetBucketLocation"
            ],
            "Resource": [
                "${aws_s3_bucket.ec2_image_builder.arn}"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "s3:Get*",
                "s3:Put*"
            ],
            "Resource": [
                "${aws_s3_bucket.ec2_image_builder.arn}/*"
            ]
        }
    ]
}
EOF
}

resource "aws_iam_instance_profile"  "sample_linux_image_builder_node" {
  name = "sample-linux-image-builder-node"
  role = aws_iam_role.sample_linux_image_builder_node.name
}

resource "aws_security_group"  "ec2_image_builder" {
  name = "ec2-image-builder-node"
  description = "For EC2 nodes used by EC2 Image Builder"
  vpc_id = var.vpc_id
  tags = {
    ManagedBy = "Terraform"
  }
}

resource "aws_security_group_rule"  "ec2_image_builder_egress_allow_all" {
  description = "Allow all outgoing traffic"
  type = "egress"
  cidr_blocks = [
    "0.0.0.0/0"
  ]
  from_port = "1"
  to_port = "65535"
  protocol = "all"
  security_group_id = aws_security_group.ec2_image_builder.id
}

resource "aws_security_group_rule"  "ec2_image_builder_ingress_allow_tcp_22_from_my_ip_address" {
  description = "Allow SSH from my IP address"
  type = "ingress"
  cidr_blocks = [
    var.my_ip_address_cidr_block
  ]
  from_port = "22"
  to_port = "22"
  protocol = "tcp"
  security_group_id = aws_security_group.ec2_image_builder.id
}

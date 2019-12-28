# About

A brief overview of what is required by EC2 Image Builder from brief experience using it.


## Non EC2 Image Builder resources required

- AMI with Amazon SSM agent installed.
- IAM instance profile and IAM role for the EC2 instance from which an AMI will be built out of.
  - This IAM role must have the following AWS Managed IAM policies attached:
    - EC2InstanceProfileForImageBuilder
    - AmazonSSMManagedInstanceCore
  - See https://docs.aws.amazon.com/imagebuilder/latest/userguide/getting-started-image-builder.html#w117aab9c11b4b3b2b3 for more details.
- S3 bucket to house logs.
- SNS topic.
- VPC and subnet where the EC2 instance will be spun up.
- Security group for the EC2 instance.
- SSH key for the EC2 instance.
- AWS License Manager configuration. It doesn't matter what rule is there, as long as the configuration is present.


## EC2 Image Builder is required

- Component. This is the code that will be executed to build the image. Use the [component.yml](/ec2-image-builder/component.yml) file.
- Image pipeline. This brings everything above together.


## Important documentation

- [Schema of Component Document](https://docs.aws.amazon.com/imagebuilder/latest/userguide/managing-image-builder-console.html#image-builder-application-documents )
- [Supported Action modules](https://docs.aws.amazon.com/imagebuilder/latest/userguide/image-builder-action-modules.html)

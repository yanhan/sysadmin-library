# About

Run Ansible playbook on EC2 instances using AWS SSM Run Command.

We will be hosting the Ansible playbook on an S3 bucket and sending logs to CloudWatch Logs.

**NOTE:** If you need to debug any of the parameters, please take a look at the code at: https://ap-southeast-1.console.aws.amazon.com/systems-manager/documents/AWS-ApplyAnsiblePlaybooks/content?region=ap-southeast-1


## Setup

- Create an S3 bucket.
- Create an IAM role for EC2 with the following managed policies attached:
  - AmazonSSMManagedInstanceCore
  - CloudWatchAgentServerPolicy

Copy the `example.yml` file and everything in `roles` to the S3 bucket, preserving the directory structure.


## IAM policy on IAM role of EC2 instance

You will also need the following policy attached to the IAM role of the EC2 instance, replacing the placeholder string with the name of the S3 bucket:
```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "s3:ListBucket",
                "s3:GetBucketLocation"
            ],
            "Resource": "arn:aws:s3:::REPLACE_WITH_S3_BUCKET_NAME"
        },
        {
            "Effect": "Allow",
            "Action": "s3:GetObject",
            "Resource": "arn:aws:s3:::REPLACE_WITH_S3_BUCKET_NAME/*"
        }
    ]
}
```

## Running

In the `run` script, replace all the placeholder strings. There are a few parameters that bear some explanation:

- `SourceInfo`: the `path` must be a HTTPS path to a directory in an S3 bucket, or the top level of the S3 bucket. This directory should contain the playbook or contain some directory structure that contain the playbook.
- `PlaybookFile`: the relative path to the Ansible playbook with respect to the directory in `SourceInfo`.
- `ExtraVariables`: if there are multiple variables, they must be separated by a single space.


## References

- https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-state-manager-ansible.html
- https://ap-southeast-1.console.aws.amazon.com/systems-manager/documents/AWS-ApplyAnsiblePlaybooks/content?region=ap-southeast-1

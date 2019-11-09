# About

Send CloudWatch Logs to S3 bucket using Firehose.


## S3 and KMS setup

Create an S3 bucket for storing the logs.

Decide if you want to use a KMS key to encrypt the contents of the bucket.


## IAM role for Firehose

Then create an IAM role with:

- The trust policy in `sts-firehose-to-s3.json`.
- The IAM policy `firehose-to-s3.json` after replacing the placeholders for the S3 bucket name and KMS key ARN.

Then create a Kinesis Firehose delivery stream and use the above IAM role for the delivery stream. If using the AWS web console, you can remove the inline policy from the "one click" setup.

You might want to set a buffer time of 60s for quicker debugging.


## IAM role for CloudWatch Logs

Now, we create an IAM role for CloudWatch Logs. If using the AWS web console, just select Lambda as the service; we can edit it later. Do not select EC2 because that will create an IAM instance profile instead of an IAM role.

Take note of the name of the IAM role. We will assume it is named `cw-logs-to-firehose`.

The IAM role should have the following:

- The trust policy in `sts-cloudwatch-logs-to-firehose.json`, replacing the placeholder for the AWS region.
- The IAM policy in `cloudwatch-logs-to-firehose.json`, replacing the placeholders for AWS region, AWS account id and IAM role name (which we assume is called `cw-logs-to-firehose`).


## CloudWatch Logs subscription filter

Using the awscli tool, run the following after replacing all the placeholder values:
```
aws logs put-subscription-filter \
	--logs-group-name REPLACE_WITH_CLOUDWATCH_LOG_GROUP_NAME \
	--filter-name myfilter \
	--filter-pattern "" \
	--destination-arn REPLACE_WITH_FIREHOSE_DELIVERY_STREAM_ARN \
	--role-arn REPLACE_WITH_IAM_ROLE_FOR_CLOUDWATCH_LOGS_TO_FIREHOSE
```

To generate logs, either go to the Firehose delivery stream and test with demo data, or wait for logs to go into the CloudWatch Log group.


## References

- https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs//SubscriptionFilters.html#FirehoseExample
- https://aws.amazon.com/premiumsupport/knowledge-center/s3-large-file-encryption-kms-key/

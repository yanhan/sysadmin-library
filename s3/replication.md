If you use Terraform to setup S3 bucket replication with KMS encryption on both source and destination buckets, you have to go the AWS web console and under replication, deselect the KMS keys, save the settings. Then edit the replication rule and select the KMS keys again.

Otherwise replication will fail.

# About

Fluent Bit DaemonSet for sending logs to CloudWatch Logs.

fluent-bit.yaml is from: https://github.com/aws-samples/amazon-cloudwatch-container-insights/blob/k8s/1.3.10/k8s-deployment-manifest-templates/deployment-mode/daemonset/container-insights-monitoring/fluent-bit/fluent-bit.yaml

create-cm.sh is from: https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/Container-Insights-setup-logs-FluentBit.html


## Tested using

- EKS; v1.22.11 (27 Aug 2022)


## IAM role for service account

Ensure that you have an IAM role for service account with the following IAM permissions:
```
{
  "Effect": "Allow"
  "Actions": [
    "logs:CreateLogStream",
    "logs:CreateLogGroup",
    "logs:PutLogEvents",
  ],
  "Resources": [
    "arn:aws:logs:AWS_REGION:AWS_ACCOUNT_ID:log-group:/aws/containerinsights/EKS_CLUSTER_NAME/*",
  ]
}
```

replacing `AWS_REGION`, `AWS_ACCOUNT_ID` and `EKS_CLUSTER_NAME` with appropriate values.


## How to use

First, create the namespace:
```
kubectl apply -f ./namespace.yaml
```

Then, open the `create-cm.sh` script with a text editor and modify the values of the `ClusterName` and `RegionName` variables. Then run it:
```
./create-cm.sh
```

Finally, modify the value of the annotation `eks.amazonaws.com/role-arn` in `fluent-bit.yaml`. Then run:
```
kubectl apply -f ./fluent-bit.yaml
```


## Viewing Logs

You can go to CloudWatch Logs Insights to query for logs. You will need to select the correct log group (it has the prefix `/aws/containerinsights/EKS_CLUSTER_NAME`); the default query should show some application logs.


## Copyright

The following files in this directory are Copyright to Amazon Web Services under the MIT-0 License:

- create-cm.sh
- fluent-bit.yaml
- namespace.yaml

All other files in this directory are Copyright to Yan Han Pang, under the 3-Clause BSD License.


## References

- https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/Container-Insights-setup-logs-FluentBit.html
- https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/Container-Insights-setup-EKS-quickstart.html

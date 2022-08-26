# About

cluster-autoscaler on EKS using ASG autodiscovery.

Repo: https://github.com/kubernetes/autoscaler

Manifests: https://raw.githubusercontent.com/kubernetes/autoscaler/cluster-autoscaler-1.23.1/cluster-autoscaler/cloudprovider/aws/examples/cluster-autoscaler-autodiscover.yaml


## Tested using

- EKS; k8s v1.22.12 (26 Aug 2022)
  - Using 1 multi-AZ ASG with mixed instances, spot lifecycle


## Configuration

First, make a copy of the `cluster-autoscaler-autodiscover.yaml` file.

### IAM policy for IAM for service account

Make a copy of the `policy.json.tmpl` file. Replace the `REGION`, `AWS_ACCOUNT_ID` and `ASG_NAME` fields accordingly.

**Pro tip:** it is probably better to use an ASG prefix instead of name. In other words, in the IAM policy, instead of using `...:autoScalingGroupName/actual-name`, use `...:autoScalingGroupName/asg-prefix*`.

Create an IAM role for service account with this policy attached. Refer to https://docs.aws.amazon.com/eks/latest/userguide/associate-service-account-role.html for details on creating this IAM role, especially wrt to its trust policy.

### EKS Service account annotation

In your copy of the `cluster-autoscaler-autodiscover.yaml` file, locate the ServiceAccount and replace the value of the `eks.amazonaws.com/role-arn` annotation with the ARN of the IAM role for the cluster-autoscaler.

### ASG tag for name

In your copy of the `cluster-autoscaler-autodiscover.yaml` file, locate the Deployment and replace `<YOUR CLUSTER NAME>` with your EKS cluster's actual name.


## How to use

Follow the instructions in the `Configuration` section.

Then run:
```
kubectl apply -f ./cluster-autoscaler-autodiscover.yaml
```


## Copyright

The following files in this directory are Copyright (c) to Cloud Native Computing Foundation under the Apache License 2.0:

- cluster-autoscaler-autodiscover.yml

All other files are Copyright (c) 2019 to Yan Han Pang, under the 3-Clause BSD License.


## References

- https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/cloudprovider/aws/README.md
- https://docs.aws.amazon.com/eks/latest/userguide/associate-service-account-role.html

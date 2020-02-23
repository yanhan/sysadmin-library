# About

Repo: https://github.com/kubernetes/autoscaler

Manifests: https://github.com/kubernetes/autoscaler/blob/106834360b5f992bfe39c96251aa5575a09d1b65/cluster-autoscaler/cloudprovider/aws/examples/cluster-autoscaler-autodiscover.yaml


## IAM permissions on EKS nodes IAM role

From https://github.com/kubernetes/autoscaler/tree/master/cluster-autoscaler/cloudprovider/aws#attach-iam-policy-to-nodegroup

Please ensure that the IAM role on the EKS nodes has the following permissions:
```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "autoscaling:DescribeAutoScalingGroups",
                "autoscaling:DescribeAutoScalingInstances",
                "autoscaling:DescribeTags",
                "autoscaling:DescribeLaunchConfigurations",
                "autoscaling:SetDesiredCapacity",
                "autoscaling:TerminateInstanceInAutoScalingGroup",
                "ec2:DescribeLaunchTemplateVersions"
            ],
            "Resource": "*"
        }
    ]
}
```

## Applying

Please modify the following line in `cluster-autoscaler-autodiscover.yml` by replacing `<YOUR CLUSTER NAME>` with your EKS cluster's actual name:
```
            - --node-group-auto-discovery=asg:tag=k8s.io/cluster-autoscaler/enabled,k8s.io/cluster-autoscaler/<YOUR CLUSTER NAME>
```


## Copyright

The following files in this directory are Copyright (c) to Cloud Native Computing Foundation under the Apache License 2.0:

- cluster-autoscaler-autodiscover.yml

All other files are Copyright (c) 2019 to Yan Han Pang, under the 3-Clause BSD License.

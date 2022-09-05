# About

XRay DaemonSet for EKS

The code is a mix of:

- https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/deploy_servicelens_CloudWatch_agent_deploy_EKS.html
- https://eksworkshop.com/intermediate/245_x-ray/daemonset.files/xray-k8s-daemonset.yaml


## Tested using

- EKS; k8s v1.22.12 (27 Aug 2022)


## How to use

### IAM role for service account

Create an IAM role for service account with the following managed IAM policy attached: `arn:aws:iam::aws:policy/AWSXRayDaemonWriteAccess` .

Take note of the value of the IAM role ARN.

### Applying

```
cp patches/serviceaccount.yaml{.example,}
# Modify patches/serviceaccount.yaml
# Specifically, the value of the ServiceAccount metadata.annotations 
# `eks.amazonaws.com/role-arn` . Its value should be the IAM role ARN
# obtained in the previous step

kubectl kustomize patches | kubectl apply -f -
```


## Note on EKS workshop

Do note that if you are trying this with the sample frontend and backend apps on EKS workshop, that everything has to be in the `default` namespace.


## Copyright

The following files in this directory are Copyright (c) 2018 Amazon.com, Inc. or its affiliates, under the Apache 2.0 License:

- base/xray.yaml

All other files in this directory are Copyright (c) 2019 to Yan Han Pang, under the 3-Clause BSD License.


## References

- https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/deploy_servicelens_CloudWatch_agent_deploy_EKS.html
- https://www.eksworkshop.com/intermediate/245_x-ray/

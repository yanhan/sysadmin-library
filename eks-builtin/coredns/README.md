# About

Guide: https://docs.aws.amazon.com/eks/latest/userguide/coredns.html


## On the ConfigMap

This applies if you are retrieving the ConfigMap freshly from EKS; we have already done it for configmap.yml in this repo.

Please replace this line:
```
proxy . /etc/resolve.conf
```

with:
```
forward . /etc/resolve.conf
```

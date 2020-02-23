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


## Copyright

The following files in this directory are Copyright (c) to Amazon Web Services:

- clusterrolebinding.yml
- clusterrole.yml
- configmap.yml
- deployment.yml
- serviceaccount.yml

All other files in this directory are Copyright (c) 2019 to Yan Han Pang, under the 3-Clause BSD License.

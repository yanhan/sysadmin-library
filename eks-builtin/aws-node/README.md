# About

EKS built-in networking plugin.

Repo: https://github.com/aws/amazon-vpc-cni-k8s

Code for manifests can be obtained from: https://github.com/aws/amazon-vpc-cni-k8s/blob/b0286f121f6996fbdc06f701c31b908bc214c2fa/config/v1.5/aws-k8s-cni.yaml

We obtained the contents of the following files from EKS:

- aws-node.yml
- clusterrolebinding.yml
- clusterrole.yml
- crd.yml
- serviceaccount.yml

We authored the following file to use aws-node without the EKS privileged PSP:

- psp.yml


## Copyright

The following files in this directory are Copyright (c) to Amazon Web Services, under Apache License 2.0:

- aws-node.yml
- clusterrolebinding.yml
- clusterrole.yml
- crd.yml
- serviceaccount.yml

All other files in this directory are Copyright (c) 2019 to Yan Han Pang, under the 3-Clause BSD License.

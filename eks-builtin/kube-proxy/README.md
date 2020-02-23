# About

**NOTE:** We are not backing up the kube-proxy configmap, because it contains the URL to the API server, which changes each time we destroy and start the EKS control plane.


## Copyright

The following files in this directory are Copyright (c) to Amazon Web Services:

- clusterrolebinding.yml
- clusterrole.yml
- configmap.yml
- daemonset.yml
- serviceaccount.yml

All other files in this directory are Copyright (c) 2019 to Yan Han Pang, under the 3-Clause BSD License.

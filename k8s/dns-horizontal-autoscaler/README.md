# About

DNS horizontal autoscaler


## Tested on

- EKS; k8s v1.22.12 (25 Aug 2022)

Manifest obtained from: https://kubernetes.io/docs/tasks/administer-cluster/dns-horizontal-autoscaling/


## How to use

Check the value of the `--target` flag in the Deployment. If the name of the coredns deployment is named something else, please change that.

```
kubectl apply -f ./manifest.yaml
```

### Tuning coredns scaling parameters

Before the first creation, adjust the `--default-params`

After the first creation, there will be a ConfigMap named `kube-dns-autoscaler` created. Adjust that instead.


## Copyright

The following files in this directory are Copyright (c) to The Kubernetes Authors under the CC BY 4.0 License:

- manifest.yaml

All other files in this directory are Copright (c) 2019 Yan Han Pang, under the 3-Clause BSD License.


## References

- https://kubernetes.io/docs/tasks/administer-cluster/dns-horizontal-autoscaling/

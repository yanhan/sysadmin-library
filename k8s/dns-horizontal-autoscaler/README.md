# About

DNS horizontal autoscaler


## Tested on

- kind; k8s v1.23.0 (25 Aug 2022)
- EKS; k8s v1.22.12 (25 Aug 2022)

Manifest obtained from: https://kubernetes.io/docs/tasks/administer-cluster/dns-horizontal-autoscaling/


## How to use

Check the value of the `--target` flag in the Deployment. If the name of the coredns deployment is named something else, please change that.

```
kubectl kustomize base | kubectl apply -f -
```

### Tuning coredns scaling parameters

The default configuration uses linear mode scaling, which creates a number of replicas proportional to the number of CPU cores and nodes in the k8s cluster.

There is also the ladder scaling mode. Details of these 2 modes and their configuration: https://github.com/kubernetes-sigs/cluster-proportional-autoscaler

Before the first creation, adjust the value of the `--default-params` in the Deployment PodTemplateSpec.

After the first creation, there will be a ConfigMap named `kube-dns-autoscaler` created. Adjust that instead.


## Copyright

The following files in this directory are Copyright (c) to The Kubernetes Authors under the CC BY 4.0 License:

- base/manifest.yaml

All other files in this directory are Copright (c) 2019 Yan Han Pang, under the 3-Clause BSD License.


## References

- https://kubernetes.io/docs/tasks/administer-cluster/dns-horizontal-autoscaling/
- https://github.com/kubernetes-sigs/cluster-proportional-autoscaler

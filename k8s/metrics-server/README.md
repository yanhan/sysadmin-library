# About

metrics-server manifests, ready for use. 2 instance HA deployment, with PDB set to 1 max unavailable.

Repo: https://github.com/kubernetes-sigs/metrics-server


## Tested using

Version: v0.6.1

- kind; Kubernetes 1.23.0 (25 Aug 2022)
- EKS; Kubernetes 1.22.12 (25 Aug 2022), 1.23.7 (05 Sep 2022)

Manifests: https://github.com/kubernetes-sigs/metrics-server/releases/download/v0.6.1/high-availability.yaml

For changes we made, search for the string `(MODIFIED)`


## Requirements

At least 2 nodes required due to pod anti affinity of `requiredDuringSchedulingIgnoredDuringExecution` in metrics-server Deployment template podspec.


## How to use

Without customization:
```
kubectl kustomize base | kubectl apply -f -
```

To customize deployment details (eg. resources requests and limits; see the section below on why you might want to do that):
```
cp patches/deployment.yaml{.example,}
# Modify patches/deployment.yaml
kubectl kustomize patches | kubectl apply -f -
```


## How to scale

Based on the "Scaling" section of: https://kubernetes-sigs.github.io/metrics-server/

For clusters up to 100 nodes, metrics-server requires 100m CPU and 200MiB memory. These should work or up to 100 nodes, 70 pods per node, 100 deployments with HPAs.

For larger clusters, recommendation is: 1m CPU and 2MiB memory per node in the k8s cluster. This is probably on a per metrics-server container basis. So if the k8s cluster has 5000 nodes, each metrics-server container should required 5 CPU and 10000MiB memory.


## Copyright

The following files in this directory are Copyright (c) to Cloud Native Computing Foundation under Apache License 2.0:

- base/high-availability.yaml

All other files in this directory are Copyright (c) 2019 to Yan Han Pang, under the 3-Clause BSD License.


## References

- https://kubernetes-sigs.github.io/metrics-server/
- https://github.com/kubernetes-sigs/metrics-server

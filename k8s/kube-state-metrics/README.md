# About

kube-state-metrics. Generates new metrics about objects within Kubernetes itself (eg. Pods, Deployments, DaemonSets, etc) for Prometheus monitoring.

Repo: https://github.com/kubernetes/kube-state-metrics


## Tested using

- kind; k8s v1.23.0 (25 Aug 2022)
- EKS; k8s v1.22.12 (25 Aug 2022)

Manifests: https://github.com/kubernetes/kube-state-metrics/tree/v2.5.0/examples/standard


## How to use

For the standard setup:
```
kubectl kustomize base | kubectl apply -f -
```

To customize the Deployment:
```
cp patches/deployment.yaml{.example,}
# Modify patches/deployment.yaml
kubectl kustomize patches | kubectl apply -f -
```


## Testing

Switch to `kube-system` namespace. Check that the container for the `kube-state-metrics` Deployment is running.

Then port forward from the `kube-state-metrics` service:

```
kubectl port-forward svc/kube-state-metrics 8080
```

Make a request to the `/metrics` endpoint:
```
curl -i http://127.0.0.1:8080/metrics
```


## How to scale

### Single pod

In the README, P99 latency on a 100 node cluster is about 0.906s, presumably using recommended settings of 0.1 CPU and 250MiB memory.

Increase CPU resource to reduce memory consumption, because kube-state-metrics has internal queues that will pile up and use more memory if they cannot be processed quickly enough.

### Sharding

Based on the README, there are 2 ways:

- Horizontal sharding
- Automated sharding

Horizontal sharding is more stable (in that it probably won't be removed by the maintainers), but all kube-state-metrics pods "will have the network traffic and the resource consumption for unmarshalling all objects, not just the ones they are responsible for. To optimize this further, the Kubernetes API would need to support sharded list/watch capabilities".

It seems that what happens is: all the k8s objects will be listed, then each shard decides whether the object is handled by itself. So there is an upfront cost of getting all the k8s objects and small processing / filtering done.

Automated sharding is experimental and the maintainers warn that it may be removed in the future. It uses a StatefulSet to determine the position of the each Pod. The same costs as stated above for Horizontal Sharding also comes into play, just that using StatefulSet offers a more automated deployment strategy.


## Copyright

The following files are Copyright (c) to Cloud Native Computing Foundation, under Apache License 2.0

- base/cluster-role-binding.yaml
- base/cluster-role.yaml
- base/deployment.yaml
- base/service-account.yaml
- base/service.yaml

All other files in this directory are Copyright (c) 2019 to Yan Han Pang, under the 3-Clause BSD License.


## References

- https://github.com/kubernetes/kube-state-metrics

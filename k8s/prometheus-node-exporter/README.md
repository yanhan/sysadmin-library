# About

Prometheus Node Exporter.

```
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
```


## Tested using

- kind; k8s v1.23.0 (25 Aug 2022)


## How the code is obtained

```
helm template pnx prometheus-community/prometheus-node-exporter --version 4.0.0 --set 'rbac.pspEnabled=false' >manifest.yaml
```


## Show configuration options

```
helm show values prometheus-community/prometheus-node-exporter --version 4.0.0
```

Use `helm template` to generate new manifests if you need to tweak it.


## How to use

```
kubectl apply -f ./manifest.yaml
```


## Testing

```
kubectl port-forward svc/pnx-prometheus-node-exporter 9100
```

Then run:
```
curl -i http://127.0.0.1:9100/metrics
```


## Copyright

The following files in this directory are Copyright (c) to the Prometheus Authors under Apache License 2.0:

- manifest.yaml

All other files in this directory are Copyright (c) to Yan Han Pang, under the 3-Clause BSD License.


## References

- https://github.com/prometheus-community/helm-charts/tree/main/charts/prometheus-node-exporter
- https://stackoverflow.com/a/56959821

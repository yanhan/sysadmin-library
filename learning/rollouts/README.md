# About

Learning Argo Rollouts: https://argoproj.github.io/argo-rollouts/

Tried on:
- k8s v1.23.0 (13 Nov 2022)


## Installation

Controller:
```
k create ns argo-rollouts
k apply -n argo-rollouts -f https://github.com/argoproj/argo-rollouts/releases/download/v1.3.1/install.yaml
```

kubectl argo rollouts plugin:
```
c -LO 'https://github.com/argoproj/argo-rollouts/releases/download/v1.3.1/kubectl-argo-rollouts-linux-amd64'
chmod u+x ./kubectl-argo-rollouts-linux-amd64
mv ./kubectl-argo-rollouts-linux-amd64 ~/bin/kubectl-argo-rollouts
```

Check that kubectl argocd rollouts plugin is working:
```
k argo rollouts version
```


## Basic usage

Follow this article: https://argoproj.github.io/argo-rollouts/getting-started/

It shows how to:
- Use kubectl argo rollouts to set a new image for a Rollout resource to trigger an update. Then `kubectl argo rollouts promote` is used to resume the rollout due to the indefinite pause used
- Abort a rollout


## Dashboard

```
k argo rollouts dashboard
```

Then go to http://127.0.0.1:3100


## Use with Prometheus analysis

NOTE: Study `manifest.yml` in detail if you are interested.

Install kube-prometheus:
```
g submodule update --init kube-prometheus
cd kube-prometheus
g checkout v0.11.0    # this needs to be done once
k create -f kube-prometheus/manifests/setup
until k get servicemonitors --all-namespaces; do date; sleep 1; echo ""; done
k create -f kube-prometheus/manifests/
```

Port forward and access it on http://127.0.0.1:9090:
```
kpf svc/prometheus-k8s -n monitoring 9090
```

Go to the goserver folder and build the images (they are the same image really):
```
cd goserver
docker build -t localhost:6001/yanhan/goprom-argo-rollouts:0.1 .
docker build -t localhost:6001/yanhan/goprom-argo-rollouts:0.2 .
docker push localhost:6001/yanhan/goprom-argo-rollouts:0.1
docker push localhost:6001/yanhan/goprom-argo-rollouts:0.2
```

Then deploy the manifest (it has a Rollout, AnalysisTemplate, Service and ServiceMonitor):
```
kaf ./manifest.yml
```

Monitor the Rollout (the initial rollout will be successful):
```
k argo rollouts get rollout myapp -w
```

Do a port forward and fire some requests to our app:
```
kpf svc/myapp -n default 8888:80
for i in {1..5}; do c -i http://127.0.0.1:8888 ; done
```

On Prometheus, ensure that the metrics `myapp_main_requests` and `myapp_main_requests_total` are captured. They should have the labels `service="myapp"` and `role="stable"`. The `myapp_main_requests` metric should also have the label `success="true"`.

Now, set a new image:
```
k argo rollouts set image myapp goserver=localhost:6001/yanhan/goprom-argo-rollouts:0.2
```

If you want to cause the analysis run to fail and exceed the failureLimit, do a port forward of the canary pod and enable its failure mode (will result in ~50% requests to the pod failing), then send some requests to it. Suppose the port 8889 is used for port forwarding this pod. We do:
```
c -i http://127.0.0.1:8889/fail/enable
c -i http://127.0.0.1:8889/fail/status # should show Enabled
while true; do c -i http://127.0.0.1:8889/ ; sleep 3; done
```

The `k argo rollouts get rollout myapp -w` should after some time, show failures (you will see a number of failures followed by a red cross). Alternatively, run this for details on the analysisrun:
```
kdes analysisruns
```

After exceeding the failureLimit, a rollback will be done to the original image.


## References

- https://argoproj.github.io/argo-rollouts/installation/
- https://dev.to/codefreshio/recover-automatically-from-failed-deployments-with-argo-rollouts-and-prometheus-metrics-1oj4
- https://argoproj.github.io/argo-rollouts/features/ephemeral-metadata/
- https://argoproj.github.io/argo-rollouts/features/analysis/#failure-conditions-and-failure-limit
- https://argoproj.github.io/argo-rollouts/features/specification/

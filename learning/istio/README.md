# About

Some notes taken when learning Istio.


## Installation

NOTE: This is purely for learning purposes and is not production ready.

```
c -LO 'https://github.com/istio/istio/releases/download/1.17.1/istio-1.17.1-linux-amd64.tar.gz'
c -LO 'https://github.com/istio/istio/releases/download/1.17.1/istio-1.17.1-linux-amd64.tar.gz.sha256'
```

Extract the first. Then compare sha256 sum with the second.

Modify `PATH` to contain the `bin` directory in the extracted tarball. Suppose that directory is `${HOME}/code/istio-1.17.1`. Then do equivalent of:
```
PATH="${PATH}:${HOME}/code/istio-1.17.1/bin"
```

Source your bashrc or equivalent when done.

Ensure your k8s cluster is up. Then run:
```
istioctl install --set profile=demo -y
kubectl label ns default istio-injection=enabled
```


## Install Bookinfo sample application

```
kaf samples/bookinfo/platform/kube/bookinfo.yaml
```

Wait for the pods in default namespace to be ready:
```
k rollout status deploy/details-v1 && \
k rollout status deploy/ratings-v1 && \
k rollout status deploy/reviews-v1 && \
k rollout status deploy/reviews-v2 && \
k rollout status deploy/reviews-v3 && \
k rollout status deploy/productpage-v1
```

Then verify:
```
k exec "$(kg po -l app=ratings -o jsonpath='{.items[0].metadata.name}')" -c ratings -- curl -sS productpage:9080/productpage | grep -o "<title>.*</title>"
```

You should see:
```
<title>Simple Bookstore App</title>
```

Create Istio Ingress Gateway for the application to be open to outside traffic:
```
kaf samples/bookinfo/networking/bookinfo-gateway.yaml
```

Then:
```
istioctl analyze
```

Expected output:
```
No validation issues found when analyzing namespace: default
```

Now we expose some ports (this uses node port because we are using kind, which does not have external load balancer):
```
export INGRESS_HOST=$(kg -n istio-system po -l istio=ingressgateway -o jsonpath='{.items[0].status.hostIP}')
export INGRESS_PORT=$(kg -n istio-system svc istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name == "http2")].nodePort}')
export SECURE_INGRESS_PORT=$(kg -n istio-system svc istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name == "https")].nodePort}')
export GATEWAY_URL=$INGRESS_HOST:$INGRESS_PORT
```

Verify on a browser by going to `http://$GATEWAY_URL/productpage`

Install destination rules to allow different versions of the Bookinfo app services to be used:
```
kaf samples/bookinfo/networking/destination-rule-all.yaml
```

To verify this, refresh the Bookinfo product page a few times on the browser. Sometimes you will see the star ratings (sometimes in red color, sometimes in black color), sometimes there wil be no star ratings.


## Install addons

Such as Kialia

```
kaf samples/addons
k -n istio-system rollout status deploy/kiali
```

Access Kiali dashboard:
```
istioctl dashboard kiali
```

Go to the `Graph` page and select the `default` namespace.

Send some traffic to view the the application architecture on Kiali dashboard:
```
for i in $(seq 1 100); do curl -s -o /dev/null "http://$GATEWAY_URL/productpage"; done
```


## References

- https://istio.io/latest/docs/setup/getting-started/
- https://istio.io/latest/docs/examples/bookinfo/

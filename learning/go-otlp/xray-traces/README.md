# About

Sample code that instruments a net/http server using Go OTLP, sends the traces to ADOT Collector, which then exports to AWS X-Ray.


## Infra prereqs

- EKS cluster
- ADOT EKS add-on
- ADOT collector deployed onto the EKS cluster

For more information about the ADOT EKS add-on and how to deploy the ADOT collector for AWS X-Ray, refer to the following docs:

- https://docs.aws.amazon.com/eks/latest/userguide/opentelemetry.html
- https://docs.aws.amazon.com/eks/latest/userguide/configure-xray.html


## Building

Replace the `yanhan/go-otlp-xray:v0.2.9` with your own repo and tag:
```
docker build -t yanhan/go-otlp-xray:v0.2.9 .
docker push yanhan/go-otlp-xray:v0.2.9
```

There is a base setup in the `base` dir. But you might want to edit the namespace, names, service port, as well as the value of the `OTLP_COLLECTOR_HOST_PORT` environment variable in the Deployment spec.

If those are ok with you, deploy to EKS using:
```
kubectl kustomize base | kubectl apply -f -
```


## Copyright

Most of the code in main.go is adapted from (to "web server"-ize the Fibonacci program in the first link):

- https://opentelemetry.io/docs/instrumentation/go/getting-started/
- https://opentelemetry.io/docs/instrumentation/go/libraries/
- https://aws-otel.github.io/docs/getting-started/go-sdk/trace-manual-instr

The following files in this directory are Copyright The OpenTelemetry Authors, under the Apache 2.0 License:

- main.go

All other files in this directory are Copyright (c) 2019 to Yan Han Pang, under the 3-Clause BSD License


## References

- https://opentelemetry.io/docs/instrumentation/go/getting-started/
- https://aws-otel.github.io/docs/getting-started/go-sdk/trace-manual-instr
- https://opentelemetry.io/docs/instrumentation/go/libraries/

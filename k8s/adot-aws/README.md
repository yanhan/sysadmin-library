# About

ADOT - AWS Distribution for Open Telemetry


## Provenance

addons-otel-permissions.yaml: from https://amazon-eks.s3.amazonaws.com/docs/addons-otel-permissions.yaml

cert-manager.yaml: from https://github.com/cert-manager/cert-manager/releases/download/v1.8.2/cert-manager.yaml

collector-config-xray.yaml: from https://github.com/aws-observability/aws-otel-community/blob/2644c20c039380b13b1bde19a2836345c0cd926c/sample-configs/operator/collector-config-xray.yaml

sample-app.yaml: from https://github.com/aws-observability/aws-otel-community/blob/2644c20c039380b13b1bde19a2836345c0cd926c/sample-configs/sample-app.yaml


## Tested using

- EKS; k8s v1.22.11 (31 Aug 2022)


## How to use

### Apply prerequisite manifests

```
kubectl apply -f ./addons-otel-permissions.yaml
kubectl apply -f ./cert-manager.yaml
```

### adot EKS addon (the ADOT Operator)

Then create EKS addon for adot. This can be done using a number of IAC tools including Terraform. See https://docs.aws.amazon.com/eks/latest/userguide/adot-manage.html

### IAM role for service account

After that (or in parallel), create IAM role for service account for the ADOT Collector. For XRay, that IAM role needs to have the `arn:aws:iam::aws:policy/AWSXrayWriteOnlyAccess` managed policy attached. See https://docs.aws.amazon.com/eks/latest/userguide/adot-iam.html

Please take note of the namespace used. The ADOT Collector needs to be deployed in the same namespace as the service account.

Take note of the IAM role, then modify `sa.yaml` with the correct value of the `eks.amazonaws.com/role-arn` annotation. Then apply:
```
kubectl apply -f ./sa.yaml
```

### Install ADOT Collector

Thereafter, modify the `collector-config-xray.yaml` file, in particular the values for these fields:

- `metadata.namespace` (must be the same as the ServiceAccount in previous step)
- `spec.config` -> `exporters.awsxray.region`.
- `spec.serviceAccount` (must be the name of the ServiceAccount in the previous step)

Then apply:
```
kubectl apply -f ./collector-config-xray.yaml
```

### Sending traces to the ADOT Collector

By default, creating the ADOT Collector in the previous step exposes a k8s Service with a gRPC endpoint at port 4317 and a HTTP endpoint at port 4318.

Applications running on any namespace can send traces to the Service using either port.

For applications running on a different namespace, remember to minimally add the `NAMESPACE.svc` domain after the Service name for DNS resolution to work correctly. eg. `my-collector-xray-collector.default.svc`. Or the full DNS name: `my-collector-xray-collector.default.svc.cluster.local`.


## Verifying

Modify `sample-app.yaml`, in particular the following values:

- `namespace` for all resources
- for the `Deployment`, various env vars in `spec.template.spec.containers[0].env`, especially `AWS_REGION` and `OTEL_EXPORTER_OTLP_ENDPOINT`

This creates a Deployment exposed via a Service. Port forward the Service:
```
kubectl port-forward svc/sample-app 4567
```

Then curl:
```
curl -i http://127.0.0.1:4567/outgoing-http-call
```

This should return a response similar to:
```
{"traceId": "1-f9abb04e-596f90a9c301bd0f1203e671"}
```

Go to AWS XRay console to search for this Trace.


## Copyright

The following files in this directory are Copyright (c) Amazon.com, Inc. or its affiliates, under the Apache 2.0 License:

- addons-otel-permissions.yaml
- collector-config-xray.yaml
- sample-app.yaml

The following files in this directory are Copyright (c) The cert-manager Authors, under the Apache 2.0 License:

- cert-manager.yaml

All other files in this directory are Copyright (c) 2019 to Yan Han Pang, under the 3-Clause BSD License.

# About

Simple Cortex setup using helm chart


## Add helm repo

```
helm repo add cortex-helm https://cortexproject.github.io/cortex-helm-chart
```

We are using chart version 1.6.0


## AWS prereqs

- EKS cluster
- S3 buckets for blocks, alertmanager, ruler
- IAM role for service account for `cortex` service account in `cortex` namespace, with permissions to Delete, Get, ListBucket, Put objects on the above-mentioned S3 buckets


## Installing helm chart

Please install the gomplate tool: https://github.com/hairyhenderson/gomplate

Create a file named `aws-data.json` with the following details:
```
{
  "region": "REPLACE WITH AWS region of S3 buckets",
  "s3": {
    "blocks": "REPLACE WITH blocks S3 bucket name",
    "alertmanager": "REPLACE WITH alertmanager S3 bucket name",
    "ruler": "REPLACE WITH ruler S3 bucket name"
  },
  "cortex_iam_role_arn": "REPLACE WITH cortex IAM role ARN"
}
```

Then, generate the `values.yaml` file using gomplate:
```
gomplate -d aws=./aws-data.json -f ./values.yaml.tmpl -o values.yaml
```

Then run:
```
kubectl apply -f ./namespace.yaml
helm install cortex cortex-helm/cortex --namespace cortex --values values.yaml --version 1.6.0
```


## Pushing metrics to Cortex

On AWS, create a TCP Target Group containing all the EKS nodes. Ensure the port is the same as the NodePort of the `cortex-nginx` service.

On the EKS node group security group, add a rule to allow ingress from `0.0.0.0` to the NodePort.

Create an NLB that points to the Target Group. HTTP port 80 on the NLB will do.

Wait for the Target Group health checks to pass.

Prometheus should push to Cortex on Cortex's `/api/v1/push` endpoint. For instance, if the URL of the NLB is `my-cortex-nlb-3c225f0283de169c.elb.ap-southeast-1.amazonaws.com`, the remote write config of the Prometheus should look similar to:
```
remote_write:
  - name: cortex
    url: http://my-cortex-nlb-3c225f0283de169c.elb.ap-southeast-1.amazonaws.com/api/v1/push
```

Query Cortex on the `/prometheus/api/v1/query` endpoint. API is same as Prometheus itself.


## Grafana Data Source

The Data Source type is Prometheus. The API endpoint is `/prometheus`. Following the above example which uses an NLB, the Data Source URL you should enter on Grafana is `http://my-cortex-nlb-3c225f0283de169c.elb.ap-southeast-1.amazonaws.com/prometheus`


## References

- https://github.com/cortexproject/cortex-helm-chart/
- https://cortexproject.github.io/cortex-helm-chart/
- https://github.com/cortexproject/cortex-helm-chart/blob/master/values.yaml
- https://github.com/cortexproject/cortex-helm-chart/blob/master/templates/_helpers.tpl

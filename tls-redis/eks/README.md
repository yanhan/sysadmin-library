# About

Skeleton code for connecting to Elasticache over TLS with password auth on EKS.


## Instructions

Create a k8s Secret object named `tls-redis` in the `default` namespace, with the following keys:

- `REDIS_ENDPOINT` (this must contain the `:6379` at the end)
- `REDIS_AUTH_TOKEN`

If using my built image on Docker Hub, run:
```
kubectl apply -f ./manifest.yml
```

## Building the image

```
docker build -t TAG_OF_YOUR_CHOOSING .
```

Then modify the `image` line in `manifest.yml` and deploy it.

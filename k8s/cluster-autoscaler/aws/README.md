# About

cluster-autoscaler on EKS using ASG autodiscovery.

Repo: https://github.com/kubernetes/autoscaler

Manifests: https://raw.githubusercontent.com/kubernetes/autoscaler/cluster-autoscaler-1.23.1/cluster-autoscaler/cloudprovider/aws/examples/cluster-autoscaler-autodiscover.yaml


## Tested using

- EKS; k8s v1.22.12 (26 Aug 2022)
  - Using 1 multi-AZ ASG with mixed instances, spot lifecycle
  - Using 2 multi-AZ ASG with mixed instances, spot lifecycle, different capacity nodes for each ASG
  - Using 3 single AZ ASG with mixed instances, spot lifecycle, with PersistentVolume backed by EBS


## NOTE on use with PersistentVolume backed by EBS

EBS volumes are single AZ only. As such, if you are using PersistentVolumes backed by EBS volumes, you will want Pods to be scheduled onto the same AZ as the EBS volume.

You will want to setup:

- Multiple ASGs (Node Groups), each ASG only with subnets from 1 single AZ
- Add the `--balance-similar-node-groups` command line flag to the cluster-autoscaler Deployment spec
- PersistentVolumes must have their `nodeAffinity.required.nodeSelectorTerms` spec with `topology.kubernetes.io/zone` set to the AZ. Refer to the example code below

Example of setting nodeAffinity on PersistentVolume spec:
```
  nodeAffinity:
    required:
      nodeSelectorTerms:
        - matchExpressions:
            - key: topology.kubernetes.io/zone
              operator: In
              values:
                - ap-southeast-1a
```

It is only with this that Kubernetes will know to schedule a Pod using this PersistentVolume onto a node in the same AZ. Which will then let cluster-autoscaler take it into account and scale up the correct ASG in the desired AZ.

The cluster-autoscaler documentation could do better in explaining this.

Credits to this GitHub issue for giving enough of a hint for me to solve it: https://github.com/kubernetes/autoscaler/issues/4739#issuecomment-1094109090


## Configuration

First, make a copy of the `cluster-autoscaler-autodiscover.yaml` file.

### IAM policy for IAM for service account

Make a copy of the `policy.json.tmpl` file. Replace the `REGION`, `AWS_ACCOUNT_ID` and `ASG_NAME` fields accordingly.

**Pro tip:** it is probably better to use an ASG prefix instead of name. In other words, in the IAM policy, instead of using `...:autoScalingGroupName/actual-name`, use `...:autoScalingGroupName/asg-prefix*`.

Create an IAM role for service account with this policy attached. Refer to https://docs.aws.amazon.com/eks/latest/userguide/associate-service-account-role.html for details on creating this IAM role, especially wrt to its trust policy.


## How to use

Follow the instructions in the `Configuration` section.

Create `patches/serviceaccount.yaml` from the example file. Then modify the value of the annotation `eks.amazonaws.com/role-arn`, using the IAM role ARN for the cluster-autoscaler:
```
cp patches/serviceaccount.yaml{.example,}
# Modify eks.amazonaws.com/role-arn annotation
```

Create `patches/deployment.yaml` from the example file. Then modify the command, replacing `<YOUR CLUSTER NAME>` with the actual name of your EKS cluster:
```
cp patches/deployment.yaml{.example,}
# Modify to replace ASG tag with actual EKS cluster name
```

Then run:
```
kubectl kustomize patches | kubectl apply -f -
```


## Copyright

The following files in this directory are Copyright (c) to Cloud Native Computing Foundation under the Apache License 2.0:

- base/cluster-autoscaler-autodiscover.yml

All other files are Copyright (c) 2019 to Yan Han Pang, under the 3-Clause BSD License.


## References

- https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/cloudprovider/aws/README.md
- https://docs.aws.amazon.com/eks/latest/userguide/associate-service-account-role.html
- https://github.com/kubernetes/autoscaler/issues/4739#issuecomment-1094109090

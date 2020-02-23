# About

EKS Privileged PSP and RBAC. For backup purposes.


## How to retrieve

```
kubectl get psp eks.privileged -o yaml
kubectl get clusterrole eks:podsecuritypolicy:privileged -o yaml
kubectl get clusterrolebinding eks:podsecuritypolicy:authenticated -o yaml
```

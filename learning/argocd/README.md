# About

ArgoCD learning


## Installation

Ensure you have installed the argocd client binary (It kind of functions like kubectl). One way is to use asdf. The plugin repo is this: https://github.com/beardix/asdf-argocd.git

To install the argocd server on k8s:
```
k create ns argocd
k apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/v2.5.1/manifests/install.yaml
```

Use port forwarding to access the argocd API server:
```
kpf -n argocd svc/argocd-server 8080:443
```

Get admin password:
```
kg -n argocd secret argocd-initial-admin-password -o jsonpath='{.data.password}' | base64 -d; echo
```

Login (username is `admin`):
```
argocd login --insecure 127.0.0.1:8080
```


## Guestbook app to get started

Example app:
```
kaf ./guestbook-app.yml
```

Sync:
```
argocd app sync guestbook --insecure
```


## echo-manifests

This is an app that you can update and see autosync. It is in the echo-manifests dir (tracked as a git submodule)

Manifest of the app: echo-app.yml



## Common questions

https://argo-cd.readthedocs.io/en/stable/faq/#how-often-does-argo-cd-check-for-changes-to-my-git-or-helm-repository

https://argo-cd.readthedocs.io/en/stable/user-guide/app_deletion/


## Read progress

Finished reading the following / skipped reading:

- Overview
- Understand The Basics (skipped; these are prereq knowledge for k8s which we know)
- Core Concepts
- Getting Started
- Developer Guide (skipped; we are not ArgoCD developers)
- FAQ
- Security Considerations
- Support

### User Guide

- Overview
- Tools
- Kustomize
- Jsonnet
- Directory
- Build Environment
- App Deletion

### Operator Manual

- Architectural Overview

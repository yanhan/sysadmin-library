## NOTES

The Vault token issued from the Kubernetes login (done by sethvargo/vault-kubernetes-authenticator) will expire according to the `ttl` in `auth/kubernetes/role/app-for-mysql-mg-vault`. We configure the consul-template container to renew it.

The files in the shared volume (which includes the the Vault token and accessor) must have have their owner set to the uid of the consul-template container. Otherwise it cannot read or renew them.

The `fsGroup` of the Pod should be set to the gid of the main container so it can read the secret files rendered by the consul-template container.

The Golang code is not thread safe. It is possible for the database connection to be used while a reconnection is in progress.


## On TTLs

The TTLs of the MySQL database role is deliberately set to a very low value (2 min) for you to witness credential rotation in the log messages.

To witness Vault token rotation by consul-template, set the TTL for the `auth/kubernetes/role/app-for-mysql-mg-vault` to a low value, say 2 min.



## To try things out locally

Minikube when using VirtualBox driver will create a network interface named `vboxnet0` with a /24 subnet.

Suppose this subnet is `192.168.77.0/24`.

Configure Consul server to listen on `192.168.77.1:8200`. Then unseal the vault from the host machine.

If using iptables, ensure that you add a rule that allows traffic from the /24 subnet to the Consul server IP address and port.

From within any container in Minikube, run:
```
curl -H 'X-Vault-Token: REPLACE_WITH_VAULT_TOKEN'  http://192.168.77.1:8200/v1/sys/mounts
```


## ConfigMap for `VAULT_ADDR`

Create a k8s Secret object called `vault-mysql-no-tls` in the `default` namespace. This Secret object has a single key `VAULT_ADDR` whose value is the Vault address.


## Install MySQL using Helm

```
helm install --name mysql-mg-vault stable/mysql
```

Expose the service to 127.0.0.1:
```
kubectl port-forward svc/mysql-mg-vault 13306:3306
```

Extract the MySQL root password from the Kubernetes Secret object.

Then login to MySQL and create a second admin user that will be managed by Vault:
```
CREATE USER 'smapper'@'%' IDENTIFIED BY 'REPLACE WITH PASSWORD';
GRANT ALL PRIVILEGES ON *.* TO 'smapper'@'%' WITH GRANT OPTION;
FLUSH PRIVILEGES;
```

Still in MySQL, create a database named `apps_galore`:
```
CREATE DATABASE apps_galore;
```

Initialize the database from the outside:
```
mysql -h 127.0.0.1 -p -P 13306 -u root -D apps_galore <./mysql-no-tls/tables.sql
```


## Let Vault administer MySQL

```
vault secrets enable database
```

Configure:
```
vault write database/config/mysql-mg-vault \
  plugin_name=mysql-database-plugin \
  connection_url="{{username}}:{{password}}@(127.0.0.1:13306)/" \
  allowed_roles="app-for-mysql-mg-vault" \
  username="REPLACE WITH second MySQL root user" \
  password="REPLACE WITH PASSWORD"
```

Rotate the password:
```
vault write -force database/rotate-root/mysql-mg-vault
```

Create a role (**NOTE:** The value of the `db_name` parameter must be the same as the `XX` in `database/config/XX`` above when we created the configuration):
```
vault write database/roles/app-for-mysql-mg-vault \
  db_name=mysql-mg-vault \
  creation_statements="CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}'; GRANT SELECT, INSERT, UPDATE ON apps_galore.* TO '{{name}}'@'%';" \
  default_ttl="2m" \
  max_ttl="2m"
```

Test out the credentials:
```
vault read database/creds/app-for-mysql-mg-vault
```


## Let Vault use Kubernetes as authentication mechanism

Create the Service Account and ClusterRoleBinding:
```
kubectl apply -f mysql-no-tls/rbac.yml
```

Create Vault policy:
```
vault policy write app-for-mysql-mg-vault ./app-for-mysql-mg-vault.hcl
```

Extract the Service Account's JWT token, CA cert and k8s IP address:
```
export VAULT_SECRET_NAME=$(kubectl get sa vault-mysql-no-tls -o jsonpath='{.secrets[0].name}')
export SA_JWT_TOKEN=$(kubectl get secret "${VAULT_SECRET_NAME}" -o jsonpath='{.data.token}' | base64 -d)
export SA_CA_CRT=$(kubectl get secret "${VAULT_SECRET_NAME}" -o jsonpath="{.data['ca\.crt']}" | base64 -d)
export K8S_HOST=$(minikube ip)
vault auth enable kubernetes
vault write auth/kubernetes/config \
  token_reviewer_jwt="${SA_JWT_TOKEN}" \
  kubernetes_host="https://${K8S_HOST}:8443" \
  kubernetes_ca_cert="${SA_CA_CRT}"
vault write auth/kubernetes/role/app-for-mysql-mg-vault \
  bound_service_account_names=vault-mysql-no-tls \
  bound_service_account_namespaces=default \
  policies=app-for-mysql-mg-vault \
  ttl=24h
```

Create a temporary pod and test you can reach Kubernetes:
```
kubectl run vtest --image alpine:3.11.3 --generator=run-pod/v1 --rm -it --serviceaccount=vault-mysql-no-tls
```

When inside, run the following to test that the Kubernetes auth to Vault is working:
```
apk add curl jq
export SA_TOKEN=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)
curl -XPOST --data '{"jwt": "'"${SA_TOKEN}"'", "role": "app-for-mysql-mg-vault"}' http://192.168.77.1:8200/v1/auth/kubernetes/login | jq
```

You should get a JSON that contains a Vault token somewhere inside.


## References

- https://www.vaultproject.io/docs/secrets/databases/mysql-maria/
- https://learn.hashicorp.com/vault/developer/db-creds-rotation
- https://caylent.com/using-hashicorp-vault-on-kubernetes
- https://github.com/sethvargo/vault-kubernetes-authenticator
- https://codelabs.developers.google.com/codelabs/vault-on-gke/index.html#21

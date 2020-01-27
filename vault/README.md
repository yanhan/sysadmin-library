## To try things out locally

Minikube when using VirtualBox driver will create a network interface named `vboxnet0` with a /24 subnet.

Suppose this subnet is `192.168.77.0/24`.

Configure Consul server to listen on `192.168.77.1:8200`. Then unseal the vault from the host machine.

If using iptables, ensure that you add a rule that allows traffic from the /24 subnet to the Consul server IP address and port.

From within any container in Minikube, run:
```
curl -H 'X-Vault-Token: REPLACE_WITH_VAULT_TOKEN'  http://192.168.77.1:8200/v1/sys/mounts
```

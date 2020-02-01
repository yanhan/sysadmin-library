## Configuration on Vault

```
vault secrets enable -path=ssh-client-signer ssh
```

This will be used to sign client certificates.

Generate CA for client certs:
```
vault write ssh-client-signer/config/ca generate_signing_key=true
```

and get its public key (no authentication required):
```
vault read -field=public_key ssh-client-signer/config/ca >./trusted-user-ca-key.pem
```

Create a Vault role for signing SSH keys of clients (**NOTE**: The value of `default_user` is important, please change it accordingly, otherwise you will not be able to SSH in):
```
vault write ssh-client-signer/roles/my-role -<<EOF
{
  "allow_user_certificates": true,
  "allowed_users": "*",
  "default_extensions": [
    {
      "permit-pty": ""
    }
  ],
  "key_type": "ca",
  "default_user": "ubuntu",
  "ttl": "30m0s"
}
EOF
```


## On the client

```
ssh-keygen -t rsa -b 4096 -C 'someone@somewhere.com' -f ./someone_id_rsa
```

Get Vault to sign the public key and save the result:
```
vault write -field=signed_key ssh-client-signer/sign/my-role public_key=@./someone_id_rsa.pub >signed-client-key.pub
```


## Configuration on SSH server

Add the following line to `/etc/ssh/sshd_config`:
```
TrustedUserCAKeys /etc/ssh/trusted-user-ca-key.pem
```

Copy the public key extracted above to `/etc/ssh/trusted-user-ca-key.pem` in the SSH server. Set its permissions to 0400.

Then restart the SSH server:
```
sudo systemctl restart ssh
```


## Go in

```
ssh -i ./signed-client-key.pub -i ./someone_id_rsa.pub username@host
```

This is supposed to work without adding the client's public key to the authorized keys file.


## Host Key signing

For the client to verify the host before attempting to SSH.

Enable the host signer:
```
vault secrets enable -path=ssh-host-signer ssh
```

Configure CA:
```
vault write ssh-host-signer/config/ca generate_signing_key=true
```

Create role for signing host keys (**NOTE**: Change the `allowed_domains` accordingly):
```
vault write ssh-host-signer/roles/hostrole \
  key_type=ca \
  ttl=87600h \
  allow_host_certificates=true \
  allowed_domains="127.0.0.1" \
  allow_subdomains=true
```

Sign the server's SSH public key:
```
vault write ssh-host-signer/sign/hostrole \
  cert_type=host \
  public_key=@./ssh_host_rsa_key.pub
```

Extract the signed certificate:
```
vault write -field=signed_key ssh-host-signer/sign/hostrole \
  cert_type=host \
  public_key=@./ssh_host_rsa_key.pub > ./ssh_host_rsa_key-cert.pub
```

Copy the signed certificate to `/etc/ssh/ssh_host_rsa_key-cert.pub` on the SSH server and update its owner and permissions.

Add these lines to `/etc/ssh/sshd_config` on the SSH server:
```
HostKey /etc/ssh/ssh_host_rsa_key
HostCertificate /etc/ssh/ssh_host_rsa_key-cert.pub
```

Restart the SSH service on the SSH server.

On the SSH client, retrieve the host signing CA key from Vault. This can be obtained without any Vault token:
```
curl "${VAULT_ADDR}/v1/ssh-host-signer/public_key"
```

or just
```
vault read -field=public_key ssh-host-signer/config/ca
```

Add the resulting public key to the `~/.ssh/known_hosts` file, in such an entry:
```
@cert-authority REPLACE_WITH_ALLOWED_DOMAIN ssh-rsa AAAA...
```

Then SSH into the server.


## References

- https://www.vaultproject.io/docs/secrets/ssh/signed-ssh-certificates/

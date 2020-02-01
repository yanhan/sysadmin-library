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


## References

- https://www.vaultproject.io/docs/secrets/ssh/signed-ssh-certificates/

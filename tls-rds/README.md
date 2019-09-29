# About

Golang code that connects to RDS MySQL over TLS.

In `main.go`, the handler for the `/` endpoint will make a connection to a table in a MySQL database to extract a value.


## RDS CA Cert bundle

Download the TLS CA cert bundle from: https://s3.amazonaws.com/rds-downloads/rds-combined-ca-bundle.pem


## Build

```
./build
```

Then push the image to Docker Hub.


## Deploying

Create a k8s Secret object named `tls-rds`, with the following data fields:

- `DB_HOST`
- `DB_USER`
- `DB_PASSWORD`

Then run:
```
kubectl apply -f ./manifest.yml
```


## References

- https://github.com/Go-SQL-Driver/MySQL/
- https://godoc.org/github.com/go-sql-driver/mysql#RegisterTLSConfig
- https://godoc.org/crypto/tls#Config
- https://godoc.org/crypto/x509#CertPool
- https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_MySQL.html#MySQL.Concepts.SSLSupport

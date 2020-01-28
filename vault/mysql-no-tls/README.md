# About

Golang code that connects to RDS MySQL. No TLS.

In `main.go`, the handler for the `/` endpoint will make a connection to a table in a MySQL database to extract a value.


## Build

```
./build
```

Then push the image to Docker Hub.


## Deploying

Create a k8s Secret object named `vault-mysql-no-tls`, with the following data fields:

- `DB_HOST`
- `DB_USER`
- `DB_PASSWORD`

Then run:
```
kubectl apply -f ./manifest.yml
```


## References

- https://github.com/Go-SQL-Driver/MySQL/

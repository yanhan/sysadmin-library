On the EC2 instance, install the RDS 2019 cert by running:

```
curl -LO 'https://s3.amazonaws.com/rds-downloads/rds-combined-ca-bundle.pem'
sudo mv rds-combined-ca-bundle.pem /etc/ssl/certs/rds-combined-ca-bundle.pem 
sudo chmod 444 /etc/ssl/certs/rds-combined-ca-bundle.pem 
```

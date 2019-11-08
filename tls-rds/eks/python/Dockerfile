FROM python:3.7.5-buster

WORKDIR /app
COPY main.py requirements.txt ./
RUN \
  curl -L 'https://s3.amazonaws.com/rds-downloads/rds-combined-ca-bundle.pem' -o /etc/ssl/certs/rds-combined-ca-bundle.pem && \
  chmod 444 /etc/ssl/certs/rds-combined-ca-bundle.pem && \
  pip install -r requirements.txt

ENTRYPOINT ["python", "/app/main.py"]

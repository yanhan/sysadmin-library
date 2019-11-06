FROM python:3.7.5-buster

WORKDIR /app
COPY main.py requirements.txt ./
RUN pip install -r requirements.txt && \
  cat /etc/ssl/certs/* > /etc/ssl/certs/ca-bundle.crt && \
  chmod 444 /etc/ssl/certs/ca-bundle.crt

ENTRYPOINT ["python", "/app/main.py"]

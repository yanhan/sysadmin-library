FROM node:12.13.0-alpine

WORKDIR /app
RUN \
  apk update && \
  apk add curl && \
  curl -L 'https://s3.amazonaws.com/rds-downloads/rds-combined-ca-bundle.pem' -o /etc/ssl/certs/rds-combined-ca-bundle.pem && \
  chmod 444 /etc/ssl/certs/rds-combined-ca-bundle.pem && \
  chown node:node /app

USER node
COPY index.js package.json /app/
RUN npm install
ENTRYPOINT ["node", "/app/index.js"]

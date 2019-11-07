FROM node:12.13.0-alpine

WORKDIR /app
RUN chown node:node /app

USER node
COPY index.js package.json /app/
RUN npm install
ENTRYPOINT ["node", "/app/index.js"]

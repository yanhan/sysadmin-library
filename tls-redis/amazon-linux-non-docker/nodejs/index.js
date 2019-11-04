const tls = require("tls");
const redis = require("redis");

const client = redis.createClient({
  host: process.env.REDIS_HOST,
  password: process.env.REDIS_AUTH_TOKEN,
  tls: {
    secureContext: tls.createSecureContext({
      ca: tls.rootCertificates
    })
  }
});

client.get("basket", function(err, reply) {
  if (err) {
    console.log(`get 'basket' error = ${err}`);
  } else {
    console.log(`get 'basket' result = ${reply.toString()}`);
  }
});

client.set("basket", "ball", function(err, reply) {
  if (err) {
    console.log(`set 'basket' error = ${err}`);
  } else {
    console.log(`set 'basket' result = ${reply.toString()}`);
  }
});

client.quit();

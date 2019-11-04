const tls = require("tls");
const { promisify } = require("util");

const Koa = require("koa");
const Router = require("koa-router");
const redis = require("redis");

const app = new Koa();

const client = redis.createClient({
  host: process.env.REDIS_HOST,
  password: process.env.REDIS_AUTH_TOKEN,
  tls: {
    secureContext: tls.createSecureContext({
      ca: tls.rootCertificates
    })
  }
});
const getAsync = promisify(client.get).bind(client);
const setAsync = promisify(client.set).bind(client);

const healthRouter = new Router();
healthRouter.get("/healthz", (ctx, next) => {
  ctx.body = {
    "status": "OK"
  };
  ctx.status = 200;
});

const getRouter = new Router();
getRouter.get("/redisget", async (ctx, next) => {
  try {
    const result = await getAsync("basket");
    console.log(`GET result = ${result}`);
    ctx.body = {
      "status": "You got me"
    };
    ctx.status = 200;
  } catch (err) {
    console.log(`GET error = ${err}`);
    ctx.body = {
      "status": "Internal Server Error"
    };
    ctx.status = 500;
  }
});

const setRouter = new Router();
setRouter.get("/redisset", async (ctx, next) => {
  var value = "ball"
  console.log(`ctx.query = ${JSON.stringify(ctx.query)}`);
  console.log(`ctx.request.query = ${JSON.stringify(ctx.request.query)}`);
  if (ctx.query && "value" in ctx.query) {
    value = ctx.query.value;
  }
  console.log(`SETTING value = ${value}`);
  try {
    const result = await setAsync("basket", value);
    console.log(`SET result = ${result}`);
    ctx.body = {
      "status": "Mutation"
    };
    ctx.status = 200;
  } catch (err) {
    console.log(`SET error = ${err}`);
    ctx.body = {
      "status": "Internal Server Error"
    };
    ctx.status = 500;
  }
});

app.use(healthRouter.routes());
app.use(healthRouter.allowedMethods());
app.use(getRouter.routes());
app.use(getRouter.allowedMethods());
app.use(setRouter.routes());
app.use(setRouter.allowedMethods());

app.listen(process.env.PORT || 8080);

app.on("exit", function() {
  client.quit();
});

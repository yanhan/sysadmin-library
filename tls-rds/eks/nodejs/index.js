const fs = require("fs");
const tls = require("tls");
const { promisify } = require("util");

const Koa = require("koa");
const Router = require("koa-router");
const mysql = require("mysql");

const app = new Koa();

const conn = mysql.createConnection({
  host: process.env.DB_HOST,
  user: process.env.DB_USER,
  password: process.env.DB_PASSWORD,
  database: process.env.DB_NAME,
  ssl: {
    ca: fs.readFileSync(process.env.RDS_TLS_CA_CERT_PATH),
  },
});

const connQuery = promisify(conn.query).bind(conn);

const healthRouter = new Router();
healthRouter.get("/healthz", (ctx, next) => {
  ctx.body = {
    "status": "OK"
  };
  ctx.status = 200;
});

const getRouter = new Router();
getRouter.get("/get", async (ctx, next) => {
  try {
    let [result, fields] = await connQuery('SELECT title FROM books LIMIT 1;');
    if (result) {
      console.log(`result = ${JSON.stringify(result)}`);
      console.log(`fields = ${fields}`);
      ctx.body = {
        "status": "Got it"
      };
      ctx.status = 200;
    } else {
      ctx.body = {
        "status": "Not found"
      };
      ctx.status = 404;
    }
  } catch (err) {
    console.log(`GET error = ${err}`);
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

app.listen(process.env.PORT || 8080);

app.on("exit", function() {
  conn.end(function(err) {
    if (err) {
      console.log(`connection end error: ${err}`);
    }
  });
});

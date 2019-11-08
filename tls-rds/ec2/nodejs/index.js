const fs = require("fs");
const { promisify } = require("util");

const mysql = require("mysql");

const rdsTlsCaCertPath = "/etc/ssl/certs/rds-combined-ca-bundle.pem"

const conn = mysql.createConnection({
  host: process.env.DB_HOST,
  user: process.env.DB_USER,
  password: process.env.DB_PASSWORD,
  database: process.env.DB_NAME,
  ssl: {
    ca: fs.readFileSync(rdsTlsCaCertPath),
  },
});

const connQuery = promisify(conn.query).bind(conn);

async function main() {
  try {
    let [result, fields] = await connQuery('SELECT title FROM books LIMIT 1;');
    if (result) {
      console.log(`result = ${JSON.stringify(result)}`);
      console.log(`fields = ${fields}`);
    } else {
      console.log("No results");
    }
  } catch (err) {
    console.log(`Error = ${err}`);
  }
  conn.end(function(err) {
    if (err) {
      console.log(`connection end error: ${err}`);
    }
  });
}

main();

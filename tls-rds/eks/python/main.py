import os
import sys

import pymysql

from flask import Flask, jsonify, request

app = Flask(__name__)

_MYSQL_CLIENT = None

@app.route("/")
def hello():
    return "Hi hi"

@app.route("/healthz")
def healthz():
    return jsonify({
        "status": "healthy",
    })

@app.route("/get")
def getmysql():
    global _MYSQL_CLIENT
    with _MYSQL_CLIENT.cursor() as cursor:
        sql_query = "SELECT title FROM books LIMIT 1;"
        cursor.execute(sql_query)
        result = cursor.fetchone()
        print("Result = {}".format(result), flush=True)
        if result is not None:
            return jsonify({
                "value": "OK",
            })
        else:
            return jsonify({
                "value": "Did not find anything",
            })

if __name__ == "__main__":
    db_host = os.environ.get("DB_HOST", None)
    db_user = os.environ.get("DB_USER", None)
    db_password = os.environ.get("DB_PASSWORD", None)
    db_name = os.environ.get("DB_NAME", None)
    rds_tls_ca_cert_path = os.environ.get("RDS_TLS_CA_CERT_PATH", None)
    must_exit = False
    if db_host is None:
        must_exit = True
        print("Env var 'DB_HOST' is missing.")
    if db_user is None:
        must_exit = True
        print("Env var 'DB_USER' is missing.")
    if db_password is None:
        must_exit = True
        print("Env var 'DB_PASSWORD' is missing.")
    if db_name is None:
        must_exit = True
        print("Env var 'DB_NAME' is missing.")
    if rds_tls_ca_cert_path is None:
        must_exit = True
        print("Env var 'RDS_TLS_CA_CERT_PATH' is missing.")
    if must_exit:
        print("Exiting.")
        sys.exit(1)
    _MYSQL_CLIENT = pymysql.connect(
        host=db_host,
        user=db_user,
        password=db_password,
        db=db_name,
        ssl={
            "ca": rds_tls_ca_cert_path,
        },
    )
    app.run(host="0.0.0.0", port=8080)

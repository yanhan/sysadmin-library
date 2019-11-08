import os
import sys

import pymysql

_RDS_TLS_CA_CERT_PATH = "/etc/ssl/certs/rds-combined-ca-bundle.pem"

def _main():
    global _RDS_TLS_CA_CERT_PATH
    db_host = os.environ.get("DB_HOST", None)
    db_user = os.environ.get("DB_USER", None)
    db_password = os.environ.get("DB_PASSWORD", None)
    db_name = os.environ.get("DB_NAME", None)
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
    if must_exit:
        print("Exiting.")
        sys.exit(1)
    mysql_client = pymysql.connect(
        host=db_host,
        user=db_user,
        password=db_password,
        db=db_name,
        ssl={
            "ca": _RDS_TLS_CA_CERT_PATH,
        },
    )
    try:
        with mysql_client.cursor() as cursor:
            sql_query = "SELECT title FROM books LIMIT 1;"
            cursor.execute(sql_query)
            result = cursor.fetchone()
            print("Result = {}".format(result), flush=True)
    finally:
        mysql_client.close()

if __name__ == "__main__":
    _main()

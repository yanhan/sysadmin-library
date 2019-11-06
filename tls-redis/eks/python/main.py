import os
import ssl

import redis

from flask import Flask, jsonify, request

app = Flask(__name__)

_REDIS_CLIENT = None

@app.route("/")
def hello():
    return "Hi hi"

@app.route("/healthz")
def healthz():
    return jsonify({
        "status": "healthy",
    })

@app.route("/getredis")
def getredis():
    global _REDIS_CLIENT
    value = _REDIS_CLIENT.get("maize")
    if value is not None:
        value = value.decode("utf-8")
    return jsonify({
        "value": value,
    })

@app.route("/setredis")
def setredis():
    global _REDIS_CLIENT
    value = request.args.get("value")
    if value is None:
        value = "herder"
    _REDIS_CLIENT.set("maize", value)
    return jsonify({
        "status": "Set value",
    })

if __name__ == "__main__":
    redis_host = os.environ.get("REDIS_HOST", None)
    redis_auth_token = os.environ.get("REDIS_AUTH_TOKEN", None)
    must_exit = False
    if redis_host is None:
        must_exit = True
        print("Env var 'REDIS_HOST' is missing.")
    if redis_auth_token is None:
        must_exit = True
        print("Env var 'REDIS_AUTH_TOKEN' is missing.")
    default_verify_paths = ssl.get_default_verify_paths()
    ca_path = default_verify_paths.cafile if default_verify_paths.cafile is not None else default_verify_paths.openssl_capath
    _REDIS_CLIENT = redis.Redis(
        host=redis_host,
        port=6379,
        password=redis_auth_token,
        ssl=True,
        # NOTE: This should be a bundle of all the CA certs.
        # The advantage of using a volume mount with the cert on EKS AMI
        # is that we can just use the AWS CA cert bundle.
        ssl_ca_certs="/etc/ssl/certs/ca-bundle.crt",
    )
    app.run(host="0.0.0.0", port=8080)

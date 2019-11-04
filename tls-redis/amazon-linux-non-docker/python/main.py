import os
import ssl

import redis

def _main():
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
    r = redis.Redis(
        host=redis_host,
        port=6379,
        password=redis_auth_token,
        ssl=True,
        ssl_ca_certs=ca_path,
    )
    r.set("maize", "herder")
    print(r.get("maize"))

if __name__ == "__main__":
    _main()

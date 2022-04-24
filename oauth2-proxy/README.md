# About

OAuth2 Proxy basics.

Tested using v7.2.1 on 24 Apr 2022

Follow instructions in the following URL to install the binary: https://oauth2-proxy.github.io/oauth2-proxy/docs/

The general instructions are similar too, just that more details are provided here.


## Getting a public domain when testing locally

The easiest way is to use ngrok.

Go to https://ngrok.com, register an account and download the binary.

Once the setup is run, run:
```
ngrok http 4180
```

This gives a public URL that we can use. Any traffic going to that URL will be redirected to your localhost's port 4180, which is where we will configure OAuth Proxy to listen on.


## Redirect URL

We will be following this format: `YOURDOMAIN/oauth2/callback`

If using ngrok and it returns `https://am86-13-47-82-115.ngrok.io`, the redirect URL will be `https://am86-13-47-82-115.ngrok.io/oauth2/callback`

This will be specified in both the `oauth2-proxy.cfg` file as well as on GitHub as the authorization callback URL (or equivalent for other OAuth providers)


## GitHub OAuth app

**NOTE:** Google OAuth is currently buggy in that it allows non test users to access an app. See:
- https://stackoverflow.com/q/66764890
- https://issuetracker.google.com/issues/211370835

This is why we choose to use the GitHub OAuth provider, which only correctly allows specific users.

For OAuth provider: use GitHub.

https://github.com/settings/developers

Go to OAuth Apps -> New OAuth App

- Application name: anything. As long as you can identify it
- Homepage URL: the ngrok URL
- Authorization callback URL: the ngrok URL followed by `/oauth2/callback`. This has to be the same value as `redirect_url` in `oauth2-proxy.cfg`
- Enable Device Flow: don't need to check this

Once created, go to the app and generate a Client secret. Take note of its value, you need this for the `oauth2-proxy.cfg` file

### GitHub users to allow

Modify the `github_users` list. Each element should be a GitHub username to allow access


## Exclude routes from OAuth

Sometimes, we want certain routes to be excluded from auth.

Supply such routes using the `skip_auth_routes` config. Possible to specify HTTP method too.


## Other settings in config

Please go through the entirety of `oauth2-proxy.cfg.example` at least once and modify the necessary values. Instructions are given there.


## Running

```
./oauth2-proxy --config=./oauth2-proxy.cfg
```


## References

- https://oauth2-proxy.github.io/oauth2-proxy/docs/
- https://stackoverflow.com/questions/10456174/oauth-how-to-test-with-local-urls

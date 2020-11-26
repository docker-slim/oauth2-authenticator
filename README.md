# OAuth2 Authenticator

Local OAuth2 authenticator application. It'll print the generated access token when the auth flow is done. The access token is also displayed in the browser the app opens. If you don't want the app to open a browser window set the `browser` flag to `manual`.

You can also use the `authenticator` (`pkg/authenticator`) package in your own app if you need a custom local auth application that uses the generated access token.

Make sure your OAuth application configuration has the correct callback URL parameter that redirects to a local server.

## Usage

### Flags

* `verbose` - verbose output (default: false)
* `browser` - browser to open the auth URL (default: default)
* `provider` - one of the supported oauth providers (default: github)
* `proto` - server protocol - http or https (default: http)
* `address` - server address (default: localhost)
* `port` - server port (default: 3000)
* `callback` - callback handler path (default: /callback)
* `id` - OAuth2 app client ID (or use CLIENT_ID env var)
* `secret` - OAuth2 app client secret (or use CLIENT_SECRET env var)
* `scopes` - comman separated list of scopes (e.g., user,email)
* `showuser` - show user information

### Examples

Basic example to authenticate with Github:

```
oauth2-authenticator \
  -id OAUTH_APP_CLIENT_ID \
  -secret OAUTH_APP_CLIENT_SECRET \
  -scopes user,email
```

Another example:

```
oauth2-authenticator \
  -provider gitlab \
  -id OAUTH_APP_CLIENT_ID \
  -secret OAUTH_APP_CLIENT_SECRET \
  -scopes user,email \
  -port 8080 \
  -callback /oauth2/callback \
  -browser firefox \
  -showuser true
```

This example shows how to authenticate with Gitlab using a custom callback server config. It also uses a non-default web browser (firefox) to open the auth URL and it prints the authenticated user info.

## Current Auth Providers

* github
* gitlab
* bitbucket
* amazon
* google

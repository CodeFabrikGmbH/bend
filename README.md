# bend

This project is supposed to offer proxy-like functionality.

The app logs the requests which are made to this server. 
The requests can be observed at the dashboard page and can be sent to other target urls at a later date.

# Usage
The dashboard is available in the `/dashboard` path.

The configuration is available in the `/configs` path.

# Configuration
The app is configured entirely through environment variables. All are optional.

| Variable | Default | Description |
|----------|---------|-------------|
| `BEND_ADDR` | `:8080` | Address the HTTP server binds to. |
| `COOKIE_SECURE` | `true` | Sets the `Secure` flag on session cookies (only sent over HTTPS). Set to `false` for plain-HTTP setups, otherwise the browser drops the login cookies. Only relevant when Keycloak auth is enabled. |
| `KEYCLOAK_HOST` | *(empty)* | Keycloak base URL. When empty, authentication is disabled and all pages are public — intended for local/dev use. |
| `KEYCLOAK_CLIENT_ID` | *(empty)* | Keycloak client id used for login. |
| `KEYCLOAK_REALM` | *(empty)* | Keycloak realm used for login. |


# Unusable paths
Some paths are not available for tracking as they are in use for internal purposes. These are:

* /dashboard
* /configs
* /login  
* /readme  
* /api/*
* /static/*
* /favicon.ico

# Path variables
It is possible to define path variables as regexes. For example, the date in the following path 
`/api/aggregate/2021-08-15` is a path variable and needs to be passed on to the target host `https://target.host`. 
Then specify the path as `/api/aggregate/\d{4}-\d{2}-\d{2}` and the target as `https://target.host`. 
The incoming request URL path will be matched against the `^/api/aggregate/\d{4}-\d{2}-\d{2}$` regex. If it matches,
the request will be forwarded to `https://target.host/api/aggregate/2021-08-15`. Capture groups and re-ordering is 
not supported.



# REM

Add reminders using a simple API.

## Setting it up with Uberspace

### HTTPS Proxying

- Setup the .htaccess file on the root

```bash
[~]$ cat .htaccess
RewriteEngine On
RewriteCond %{HTTPS} !=on
RewriteCond %{ENV:HTTPS} !=on
RewriteRule .* https://%{SERVER_NAME}%{REQUEST_URI} [R=301,L]
RewriteRule ^r/(.*) http://localhost:42888/$1 [P]
```

### Daemontools

- Setup daemontools to use the compiled binary from this repo

```bash
[~]$ uberspace-setup-service my-r ~/bin/r
```

## Configuring reminders script

I use pushover since they have a simple API.

```bash
[~]$ cat bin/pushover
#!/usr/bin/env bash
TOKEN=<TOKEN>
USER=<USER>
MESSAGE=$1

curl -s --form-string "token=${TOKEN}" --form-string "user=${USER}" --form-string "message=${MESSAGE}" https://api.pushover.net/1/messages.json
```

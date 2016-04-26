# REM

With REM you can send yourself (or somebody else) a reminder. 

All you need to do is let REM run as a daemon on a server somewhere and just navigate to the URL with your browser and your good to go!

**Examples**

- You want to remind yourself to buy milk after work today? Just visit the following link with your browser.

```
https://cip.li/rem?time=1800&message=buy milk
```

- You want to remind [me](https://cip.li/people/stefano) to wish you a happy birthday this Saturday at 13:00 since I never check Facebook?

```
https://cip.li/rem?time=1300&day=saturday&message=Wish me a Happy Birthday!
```

Rem uses the Unix `date` command in the background, so you can use its syntax to choose a day and/or time.

## Installation

Feel free to use the precompiled binary that is supplied with the repo, or just `go build` your own instead.

## Using REM with [Uberspace](https://uberspace.de/prices)

After logging into to your uberspace account via ssh, you will have to setup Proxy-Rewrite for apache and then daemontools to manage *rem*.

### HTTPS Proxying

- Setup the .htaccess file on the document root

```bash
[~]$ cat .htaccess
RewriteEngine On
RewriteCond %{HTTPS} !=on
RewriteCond %{ENV:HTTPS} !=on
RewriteRule .* https://%{SERVER_NAME}%{REQUEST_URI} [R=301,L]
RewriteRule ^rem/(.*) http://localhost:42888/$1 [P]
```

### Daemontools

- Setup daemontools to use the compiled binary from this repo

```bash
[~]$ uberspace-setup-service my-rem ~/bin/rem
```

## Configuring the script that REM will execute

The script will send a reminder to my smartphone using the Pushover service. I chose Pushover because of their simple API.

```bash
[~]$ cat bin/rem_script
#!/usr/bin/env bash
TOKEN=<TOKEN>
USER=<USER>
MESSAGE=$1

curl -s --form-string "token=${TOKEN}" --form-string "user=${USER}" --form-string "message=${MESSAGE}" https://api.pushover.net/1/messages.json
```

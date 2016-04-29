# REM

Use REM to send reminders to yourself, or someone else.

Just run the REM daemon on your server and go to the URL with your browser or `POST` to the API instead.

Rem uses the Unix `date` command in the background, so you can use its syntax to choose a day and/or time.

## Installation

Let's assume you own the domain `cip.li` and you're using [Uberspace](https://uberspace.de) as your hosting provider. You're document root is located at `/home/user/cip.li` and you want to run REM on `https://cip.li/rem`.

### 1. Clone this github repo

#### Using Uberspace

- login via ssh and clone the repo

```bash
[user@spica ~]$ git clone https://github.com/sdaros/rem ~/cip.li/rem
```

### 2. Setup HTTP Proxying and proxy `/rem` to the REM daemon

#### Using Uberspace

- Setup the .htaccess file on the document root like so

```bash
[user@spica ~]$ cat ~/cip.li/.htaccess
RewriteEngine On
RewriteCond %{HTTPS} !=on
RewriteCond %{ENV:HTTPS} !=on
RewriteRule .* https://%{SERVER_NAME}%{REQUEST_URI} [R=301,L]
RewriteRule ^rem/(.*) http://localhost:42888/$1 [P]
```

### 3. Configure REM to use a process supervisor

#### Using Uberspace

- Supervise/Run REM using daemontools

```bash
[user@spica ~]$ uberspace-setup-service rem ~/cip.li/rem/rem
Creating the ~/etc/run-rem/run service run script
Creating the ~/etc/run-rem/log/run logging run script
Symlinking ~/etc/run-rem to ~/service/rem to start the service
Waiting for the service to start ... 1 2 3 4 5 6 started!

Congratulations - the ~/service/rem service is now ready to use!
To control your service you'll need the svc command (hint: svc = service control):
...
```

### 4. Provide a script for REM to execute

#### Using [Pushover](https://pushover.net)

- This script sends reminder to my smartphone using the [Pushover](https://pushover.net) API.

```bash
[~]$ cat ~/cip.li/rem/rem_script
#!/usr/bin/env bash
TOKEN=<TOKEN>
USER=<USER>
MESSAGE=$1

curl -s --form-string "token=${TOKEN}" --form-string "user=${USER}" --form-string "message=${MESSAGE}" https://api.pushover.net/1/messages.json
```

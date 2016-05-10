# REM

Use REM to send reminders to yourself, or someone else.

Just run the REM daemon on your server then go to the URL with your browser or `POST` to the API instead.

Notifications will be sent to your smartphone if you are using the [Pushover](http://pushover.net) service.

## Installation (Using Digitalocean)

Let's assume you own the domain `cip.li` and you're using [Uberspace](https://uberspace.de) as your hosting provider. You're user name is `bob`, the document root is located at `/home/bob/cip.li` and you want to run REM on `https://cip.li/rem`.

### 1. Create a Docker Droplet then clone the github repo

Login via ssh and clone the repo to your home directory

```bash
[root@digitalocean ~]$ git clone https://github.com/sdaros/rem ~/rem && cd ~/rem
```

### 2. Customise the config file

Customise `rem.conf.example`

```bash
[root@digitalocean ~/rem]$ cp rem.conf.example rem.conf
[root@digitalocean ~/rem]$ vim rem.conf # Configure it to suit your needs
{
	"ApiToken": "n1VrLLmRMPStaX3pA8TPdh2Kl2QS3q", # Needed for the https://pushover.net Notification Service
	"ApiUser": "cf3YtkHfnSQkYb8GTWSZuPrddTPymQ", # Needed for the https://pushover.net Notification Service
	"DocumentRoot": "/app",
        "Domain": "https://cip.li",
        "NotificationApi": "https://api.pushover.net/1/messages.json",
	"Path": "rem",
	"Port": ":42888"
}
```

### 3. Configure rem to use [caddy](https://caddyserver.com) as reverse proxy

If you do not already have a reverse proxy setup for your domain name, you can use `docker-compose` to setup caddy to proxy requests for REM

- create a new Caddyfile

```bash
[root@digitalocean ~/rem]$ vim Caddyfile
# A Caddyfile for our example could look like the following:
cip.li {
        proxy /rem rem:42888 {
                proxy_header Host {host}
                proxy_header X-Real-IP {remote}
                proxy_header X-Forwarded-Proto {scheme}
        }
        gzip
        tls email_for_lets_encrypt@cip.li
}
```
- download then run `docker-compose` with `-p <YOUR_DOMAIN>`

```bash
[root@digitalocean ~/rem]$ curl -L https://github.com/docker/compose/releases/download/1.7.0/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose && chmod +x /usr/local/bin/docker-compose
[root@digitalocean ~/rem]$ docker-compose -p cipli up -d
```

If you already have a reverse proxy setup, just build the rem docker image and start it

```bash
[root@digitalocean ~/rem]$ docker run -v ./rem.conf:/app/.config/rem/rem.conf -d -p 42888:42888 --name rem sdaros/rem
```

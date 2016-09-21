# REM

Use REM to send reminders to yourself, or someone else.

Just run the REM daemon on your server then go to the URL with your browser or `POST` to the API instead.

Notifications will be sent to your smartphone if you are using the [Pushover](http://pushover.net) service.

## Installation (Quick and Dirty)

Let's assume you own the domain `cip.li` and want to run the REM daemon on Port 42888 on your RaspberryPi 3

### 1. Download the latest release to your RaspberryPi

```bash
[user@raspberrypi3 ~]$  curl -LO https://github.com/sdaros/rem/releases/download/v0.6.0/rem-Linux-armv7l && chmod +x rem-Linux-armv7l
```

### 2. Customise the config file

Create the file `~/.config/rem/rem.conf` and configure it to suit your needs
```bash
[user@raspberrypi3 ~]$ mkdir -p ~/.config/rem && vim ~/.config/rem/rem.conf
{
	"ApiToken": "n1VrLLmRMPStaX3pA8TPdh2Kl2QS3q", # Needed for the https://pushover.net Notification Service
	"ApiUser": "cf3YtkHfnSQkYb8GTWSZuPrddTPymQ", # Needed for the https://pushover.net Notification Service
	"Domain": "cip.li",
	"NotificationApi": "https://api.pushover.net/1/messages.json",
	"Path": "",
	"Port": ":42888"
}
```

### 3. Run the binary as a service

```bash
[user@raspberrypi3 ~]$ sudo systemd-run --set-env=HOME=${HOME} --unit=rem rem-Linux-armv7l
Running as unit rem.service.
[user@raspberrypi3 ~]$
```

## Installation (Using Digitalocean)

So, now were using a [Digitalocean](https://digitalocean.com) Docker droplet and want the REM daemon to listen on `https://cip.li/rem`.

### 1. Create a Docker Droplet then clone the github repo

Login via ssh and clone the repo to your home directory

```bash
[root@digitalocean ~]$ git clone https://github.com/sdaros/rem ~/rem && cd ~/rem
```

### 2. Customise the config file

Customise `rem.conf`

```bash
[root@digitalocean ~/rem]$ vim rem.conf # Configure it to suit your needs
{
	"ApiToken": "n1VrLLmRMPStaX3pA8TPdh2Kl2QS3q", # Needed for the https://pushover.net Notification Service
	"ApiUser": "cf3YtkHfnSQkYb8GTWSZuPrddTPymQ", # Needed for the https://pushover.net Notification Service
	"Domain": "cip.li",
	"NotificationApi": "https://api.pushover.net/1/messages.json",
	"Path": "/rem",
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
[root@digitalocean ~/rem]$ curl -LO $(curl -si https://api.github.com/repos/docker/compose/releases/latest | grep -oP 'http[^"]+'"`uname -s`"-"`uname -m`") > /usr/local/bin/docker-compose && chmod +x /usr/local/bin/docker-compose
[root@digitalocean ~/rem]$ docker-compose -p cipli up -d
```

### 4. If you already have a reverse proxy setup

- run the REM docker image

```bash
[root@digitalocean ~/rem]$ docker run -v ./rem.conf:/app/.config/rem/rem.conf -d -p 42888:42888 --name rem sdaros/rem
```

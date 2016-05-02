# REM

Use REM to send reminders to yourself, or someone else.

Just run the REM daemon on your server (or [build as a docker image](https://docs.docker.com/engine/understanding-docker/)) then go to the URL with your browser or `POST` to the API instead.

Rem uses the GNU Coreutils `date` command in the background, so you can use its syntax when choosing a datetime. Notifications will be sent to your smartphone if you are using the [Pushover](http://pushover.net) service.

## Installation

Let's assume you own the domain `cip.li` and you're using [Uberspace](https://uberspace.de) as your hosting provider. You're user name is `bob`, the document root is located at `/home/bob/cip.li` and you want to run REM on `https://cip.li/rem`.

### 1. Clone this github repo then edit the `rem.conf` config file

#### Using Uberspace

- login via ssh and clone the repo to your document root

```bash
[user@spica ~]$ git clone https://github.com/sdaros/rem ~/cip.li/rem && cd ~/cip.li/rem
```

### 2. Customise the config file

- customise `rem.conf.example` then copy it into `~/.config/rem/rem.conf`.

```bash
[user@spica ~/cip.li/rem]$ mkdir -p ~/.config/rem
[user@spica ~/cip.li/rem]$ vim rem.conf.example # Configure it to suit your needs
[user@spica ~/cip.li/rem]$ cp rem.conf.example ~/.config/rem/rem.conf && cat ~/.config/rem/rem.conf
{
	"ApiToken": "n1VrLLmRMPStaX3pA8TPdh2Kl2QS3q", # Needed for the https://pushover.net Notification Service
	"ApiUser": "cf3YtkHfnSQkYb8GTWSZuPrddTPymQ", # Needed for the https://pushover.net Notification Service
	"DocumentRoot": "/home/bob/cip.li",
        "NotificationApi": "https://api.pushover.net/1/messages.json",
	"Path": "/rem",
	"Port": ":42888"
}
```

### 3. Configure HTTP Proxying

- Add the following to the .htaccess file in your Document Root (for our example `/home/bob/cip.li/.htaccess`)

```bash
[user@spica ~/cip.li/rem]$ cat /home/bob/cip.li/.htaccess
RewriteEngine On
-RewriteCond %{HTTPS} !=on
-RewriteCond %{ENV:HTTPS} !=on
-RewriteRule .* https://%{SERVER_NAME}%{REQUEST_URI} [R=301,L]
-RewriteRule ^rem/(.*) http://localhost:42888/$1 [P]
```

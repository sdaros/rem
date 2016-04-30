# REM

Use REM to send reminders to yourself, or someone else.

Just run the REM daemon on your server and go to the URL with your browser or `POST` to the API instead.

Rem uses the Unix `date` command in the background, so you can use its syntax when choosing a datetime.

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
	"ApiToken": "a1VrLLmRMPStaX3pA8TPdh2Kl2QS3q",
	"ApiUser": "cf3YtkHfnSQkYb8GTWSZuPrddTPymQ",
	"DocumentRoot": "/home/bob/cip.li",
	"Path": "/rem",
	"Port": ":42888",
	"RemScript": "/home/user/.config/rem/rem_script"
}
```

### 3. Run the REM init script if using Uberspace

- `rem -init` will print the `init_script.template` bash script to standard out using the configuration parameters provided by `~/.config/rem/rem.conf`

```bash
[user@spica ~/cip.li/rem]$ ./rem -init | bash
Creating the ~/etc/run-rem/run service run script
Creating the ~/etc/run-rem/log/run logging run script
Symlinking ~/etc/run-rem to ~/service/rem to start the service
Waiting for the service to start ... 1 2 3 4 5 6 started!

Congratulations - the ~/service/rem service is now ready to use!
To control your service you'll need the svc command (hint: svc = service control):
...
```

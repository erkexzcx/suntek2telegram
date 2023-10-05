# suntek2telegram

Blackbox FTP and SMTP server for Suntek trail cameras that forward uploaded pictures to Telegram.

# Motyvation

On AliExpress there are handful of trail cameras to pick from. The one I purchased is `Suntek HC900LTE` which has 4 modes in total:
1. Storing higher quality images to SD card. (always enabled)
2. Sending them via MMS (MMS is not cheap/free where I live - not an option)
3. SMTP (sending images via SMTP server, as attachments to emails)
4. FTP (uploading to FTP, but the **implementation is broken** and it stops working in 3-48 hours).

The idea is to get pictures from this camera in a convenient way - Telegram. Incredibly useful if you are not alone who is interested in those pictures and you can get pictures in Telegram group.

Note that this software should work with all Suntek cameras that have SIM card for network connectivity and configured with `MMSCONFIG` software:

_insert_mmsconfig_image_here_

# Getting started

## Requirements

Linux server with public IP and Docker.

Suggestions:
* If you have public IP at home - with some port forwarding and home server (such as RPI) you can easily host this application.
* Alternatively, feel free to purchase instance at your favorite cloud provider. I prefer `Linode` for their cheapest shared CPU instances.

**Note** that neither method (FTP or SMTP) actually does anything beyond this application's memory. FTP method does **not** alter server files in any way and SMTP does **not** (or attempt to) send any mails. The protocol is only implemented for camera to think it's the real service, but it's actually a blackbox and the only thing it does is simply forward picture to Telegram.

## Configuration - camera

### Preparation

FTP (at least with my camera) is broken and unusable, so I am stuck with SMTP. Basically decide which one you would like to use. Changing it is easy in this application's configuration, but might be difficult to change in camera (pull out SD card, generate config, upload config, insert SD card etc.)

`MMSCONFIG` software (used to configure camera) can be found online. Google for `MMSCONFIG` and you will find a few results from `cnsuntek.com`.

### FTP method

Here is how configuration looks like:

_insert_mmsconfig_ftp_image_here_

TODO

### SMTP method

Here is how configuration looks like:

_insert_mmsconfig_smtp_image_here_

TODO

## Configuration - suntek2telegram app

Create `config.yml` file out of `config.example.yml` file.

Then create Telegram bot and add Telegram API Key as well as ChatID.

Then, if using FTP method:
1. Ensure `ftp.enabled: true`.
2. Update `ftp.bind_host` and `ftp.bind_port` if needed.
3. Set `ftp.username` and `ftp.password` to match of what you configured in camera FTP configuration.
4. Update FTP passive ports range (if needed). Ensure it is also port-forwarded directly as is if your server is behind the NAT.

Or if using SMTP method:
1. Ensure `smtp.enabled: true`.
2. Update `smtp.bind_host` and `smtp.bind_port` if needed.
3. Set `smtp.username` and `smtp.password` to match of what you configured in camera SMTP configuration.

# Installation

Only Docker images are provided. I am too lazy to provide actual binary releases. :(

## Docker compose

```yaml
services:
  suntek2telegram:
    image: ghcr.io/erkexzcx/suntek2telegram:latest
    container_name: suntek2telegram
    restart: always
    volumes:
      - ./suntek2telegram/config.yml:/config.yml
    ports:
      # You can use any TCP port for FTP or SMTP:
      - 8123:8123/tcp
      # Port range is only used for FTP PassivePorts functionality:
      #- 4100-4199:4100-4199/tcp
    environment:
      - TZ=Europe/Vilnius
```

# Misc

## What's up with FTP support

I am a bit too lazy to provide actual logs, but here is how communication works with FTP at first:
1. Device connects to FTP server.
2. Gets FTP os info.
3. Sets current working directory to specified directory in camera's config.
4. Switches to binary mode.
5. Switches to passive mode and gets passive port from FTP server.
6. Uploads file.
7. Exits.

After 3-48 hours, the implementation breaks and it basically becomes steps 1-3 only. Looking at errors file (on sd card) it shows that it tries to upload, but getting error from modem command, which was translated by ChatGPT to be "FTP server not found". Also tried explicitly closing connection after each uploaded image and notifying FTP client that connection is closing prior that - nothing has worked. This sounds like firmware/software bug which cannot be fixed by me.

Initially this software was written for FTP functionality only, but after giving up with FTP, I decided to look into SMTP method and it worked just as good as FTP, just never stopped working. So now you know a bit of backstory of this application. :)

# suntek2telegram

- [Key features](#key-features)
- [Suntek cameras](#suntek-cameras)
  * [Configuration](#configuration)
  * [Issue with FTP implementation](#issue-with-ftp-implementation)
- [Getting started](#getting-started)
  * [Requirements](#requirements)
  * [Configuration - camera](#configuration---camera)
    + [Preparation](#preparation)
    + [FTP method](#ftp-method)
    + [SMTP method](#smtp-method)
  * [Configuration - suntek2telegram app](#configuration---suntek2telegram-app)
- [Installation](#installation)
  * [Docker compose](#docker-compose)

This application serves as a blackbox FTP and SMTP server specifically designed for Suntek trail cameras, but should work with other SIM-enabled cameras too. Its primary function is to forward the pictures uploaded by the cameras to a designated Telegram (personal or group) chat.

# Key features

* **FTP Server**: The FTP server is designed to receive images from the trail cameras. However, it doesn't alter or store any files on the server. It simply acts as a conduit to receive the images and forward them to Telegram.
* **SMTP Server**: The SMTP server is implemented to make the camera believe it's interacting with a real email service. However, it doesn't send any emails. Its sole purpose is to receive the images from the camera and forward them to Telegram.

In essence, this application acts as a blackbox, mimicking the behavior of real FTP and SMTP servers to receive images from Suntek trail cameras and forward them to Telegram.

# Suntek cameras

## Configuration

This software should work with all Suntek cameras that have SIM card for network connectivity and configured with `MMSCONFIG` software:

![preview](https://github.com/erkexzcx/suntek2telegram/blob/main/images/mmsconfig.png?raw=true)

## Issue with FTP implementation

On AliExpress there are handful of trail cameras to pick from. The one I purchased is `Suntek HC900LTE`.

I am a bit too lazy to provide actual logs, but here is how communication works with FTP at first:
1. Device connects to FTP server.
2. Gets FTP os info.
3. Sets current working directory to specified directory in camera's config.
4. Switches to binary mode.
5. Switches to passive mode and gets passive port from FTP server.
6. Uploads file.
7. Exits.

After 3-48 hours, the implementation _breaks_ and it basically becomes steps 1-3 only. Looking at errors file (on sd card) it shows that it tries to upload, but getting error from modem command, which was translated by ChatGPT to be "FTP server not found". Also tried explicitly closing connection after each uploaded image and notifying FTP client that connection is closing prior that - nothing has worked. This sounds like firmware/software bug which cannot be fixed by me. As usually, Suntek support were almost clueless of what is FTP to begin with, so reaching out to them for assistance is out of question.

Initially this software was written for FTP functionality only, but after giving up with FTP, I decided to look into SMTP method and it worked just great!

# Getting started

## Requirements

Linux server with public IP and Docker.

Suggestions:
* If you have public IP at home - with some port forwarding and home server (such as RPI) you can easily host this application.
* Alternatively, feel free to purchase instance at your favorite cloud provider. I prefer `Linode` for their cheapest shared CPU instances.

## Configuration - camera

### Preparation

FTP (at least with my camera) is broken and unusable, so I am stuck with SMTP. Basically decide which one you would like to use. Changing it is easy in this application's configuration, but might be difficult to change in camera (pull out SD card, generate config, upload config, insert SD card etc.)

`MMSCONFIG` software (used to configure camera) can be found online. Google for `MMSCONFIG` and you will find a few results from `cnsuntek.com`.

### FTP method

Here is how configuration looks like:

![preview](https://github.com/erkexzcx/suntek2telegram/blob/main/images/mmsconfig_ftp.png?raw=true)

APN is configured to your SIM provider (for mobile internet) and the rest points to your FTP server.

### SMTP method

Here is how configuration looks like:

![preview](https://github.com/erkexzcx/suntek2telegram/blob/main/images/mmsconfig_smtp.png?raw=true)

APN is configured to your SIM provider (for mobile internet) and the rest points to your SMTP server. `Email Setting` needs to contain a single random email address.

## Configuration - suntek2telegram app

Create `config.yml` file out of `config.example.yml` file.

Then create Telegram bot and add Telegram API Key as well as ChatID.

Then, if using FTP method:
1. Ensure `ftp.enabled: true`.
2. Update `ftp.bind_host` and `ftp.bind_port` if needed.
3. Set `ftp.username` and `ftp.password` to match of what you configured in camera FTP configuration.
4. If needed, update `public_ip` and `passive_ports` fields too.

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

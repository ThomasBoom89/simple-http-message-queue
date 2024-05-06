# Simple Http Message Queue

This is a simple message queue / message broker based on http

![Go](https://img.shields.io/github/go-mod/go-version/thomasboom89/simple-http-message-queue/main?style=for-the-badge)
![License](https://img.shields.io/badge/license-GNU%20GPLv3-green?style=for-the-badge)

## Setup

### Docker

Modify compose.yml depending on your needs (reverse proxy config, port mapping etc.)
Then start the docker container via

```zsh
 make live
 ```

## Usage

You can publish messages with http POST on /publish message must be sent in plaintext

To receive messages you can http GET on /subscribe or connect via websocket to /ws
On the websocket you need to send "next" in plaintext to get a new message

## License

Simple Http Message Queue \
Copyright (C) 2024 ThomasBoom89. GNU GPLv3 license.

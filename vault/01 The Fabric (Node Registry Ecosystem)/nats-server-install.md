# NATS Fundamentals Workshop, Part I

## Install the Server Manually

Download 

https://github.com/nats-io/nats-server/releases/tag/v2.14.2

Extract

Put in PATH

Verify

nats-server --version


## Install the NATS Server

Linux

	curl -L https://github.com/nats-io/nats-server/releases/download/v2.14.2/nats-server-v2.14.2-linux-amd64.tar.gz -o nats-server.tar.gz
	tar -xzf nats-server.tar.gz
	sudo mv nats-server-v2.14.2-linux-amd64/nats-server /usr/local/bin/
	nats-server --version

MacOS (Apple chip)

	curl -L https://github.com/nats-io/nats-server/releases/download/v2.14.2/nats-server-v2.14.2-darwin-arm64.tar.gz -o nats-server.tar.gz
	tar -xzf nats-server.tar.gz
	sudo mv nats-server-v2.14.2-darwin-arm64/nats-server /usr/local/bin/
	nats-server --version

Windows (PowerShell)

	Invoke-WebRequest -Uri https://github.com/nats-io/nats-server/releases/download/v2.14.2/nats-server-v2.14.2-windows-amd64.zip -OutFile nats-server.zip
	Expand-Archive nats-server.zip -DestinationPath nats-server
	move nats-server\nats-server-v2.14.2-windows-amd64\nats-server.exe C:\Windows\System32\
	nats-server --version


MacOS (Intel chip)

Change `arm64` to `amd64` in the URL

	curl -L https://github.com/nats-io/nats-server/releases/download/v2.14.2/nats-server-v2.14.2-darwin-amd64.tar.gz -o nats-server.tar.gz
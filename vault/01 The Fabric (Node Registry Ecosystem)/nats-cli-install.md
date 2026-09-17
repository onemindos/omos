# NATS Fundamentals Workshop, Part I

## Install Manually

Download the binary for your OS

https://github.com/nats-io/natscli/releases/tag/v0.4.0 


Extract it

Put it in your PATH

Verify

	nats --version


## Install the NATS CLI Using cURL

Linux

	curl -L https://github.com/nats-io/natscli/releases/download/v0.4.0/nats-0.4.0-linux-amd64.zip -o nats-cli.zip
	unzip nats-cli.zip
	sudo mv nats-0.4.0-linux-amd64/nats /usr/local/bin/
	nats --version

MacOS (Apple chip)

	curl -L https://github.com/nats-io/natscli/releases/download/v0.4.0/nats-0.4.0-darwin-arm64.zip -o nats-cli.zip
	unzip nats-cli.zip
	sudo mv nats-0.4.0-darwin-arm64/nats /usr/local/bin/
	nats --version


Windows (PowerShell)

	Invoke-WebRequest -Uri https://github.com/nats-io/natscli/releases/download/v0.4.0/nats-0.4.0-windows-amd64.zip -OutFile nats-cli.zip
	Expand-Archive nats-cli.zip -DestinationPath nats-cli
	move nats-cli\nats-0.4.0-windows-amd64\nats.exe C:\Windows\System32\
	nats --version
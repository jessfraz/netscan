# netscan

[![Travis CI](https://travis-ci.org/jessfraz/netscan.svg?branch=master)](https://travis-ci.org/jessfraz/netscan)

Scan a network for ports that are open on an ip/ip range, and
ips that are in use on that network.

## Installation

#### Binaries

- **darwin** [386](https://github.com/jessfraz/netscan/releases/download/v0.0.0/netscan-darwin-386) / [amd64](https://github.com/jessfraz/netscan/releases/download/v0.0.0/netscan-darwin-amd64)
- **freebsd** [386](https://github.com/jessfraz/netscan/releases/download/v0.0.0/netscan-freebsd-386) / [amd64](https://github.com/jessfraz/netscan/releases/download/v0.0.0/netscan-freebsd-amd64)
- **linux** [386](https://github.com/jessfraz/netscan/releases/download/v0.0.0/netscan-linux-386) / [amd64](https://github.com/jessfraz/netscan/releases/download/v0.0.0/netscan-linux-amd64) / [arm](https://github.com/jessfraz/netscan/releases/download/v0.0.0/netscan-linux-arm) / [arm64](https://github.com/jessfraz/netscan/releases/download/v0.0.0/netscan-linux-arm64)
- **solaris** [amd64](https://github.com/jessfraz/netscan/releases/download/v0.0.0/netscan-solaris-amd64)
- **windows** [386](https://github.com/jessfraz/netscan/releases/download/v0.0.0/netscan-windows-386) / [amd64](https://github.com/jessfraz/netscan/releases/download/v0.0.0/netscan-windows-amd64)

#### Via Go

```bash
$ go get github.com/jessfraz/netscan
```

## Usage

```console
$ netscan -h
NAME:
   netscan - Scan network ips and ports.

USAGE:
   netscan [global options] command [command options] [arguments...]

VERSION:
   version v0.0.0, build 30bb0f0

AUTHOR(S):
   @jessfraz <no-reply@butts.com>

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug, -d          run in debug mode
   --timeout, -t "1s"   override timeout used for check
   --port, -p "1-1000"  port range to check
   --proto "tcp,udp"    protocol/s to check
   --help, -h           show help
   --version, -v        print the version
```

Examples:

```console
# for a cidr
$ netscan 192.168.0.1/24

# for a single ip
$ netscan 192.168.104.30
```

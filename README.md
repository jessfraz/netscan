# netscan

[![Travis CI](https://img.shields.io/travis/jessfraz/netscan.svg?style=for-the-badge)](https://travis-ci.org/jessfraz/netscan)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://godoc.org/github.com/jessfraz/netscan)

Scan a network for ports that are open on an ip/ip range, and
ips that are in use on that network.

 * [Installation](README.md#installation)
      * [Binaries](README.md#binaries)
      * [Via Go](README.md#via-go)
 * [Usage](README.md#usage)
      * [Examples](README.md#examples)

## Installation

#### Binaries

For installation instructions from binaries please visit the [Releases Page](https://github.com/jessfraz/netscan/releases).

#### Via Go

```console
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
   version v0.3.1, build 30bb0f0

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

#### Examples

```console
# for a cidr
$ netscan 192.168.0.1/24

# for a single ip
$ netscan 192.168.104.30
```

# netscan

[![make-all](https://github.com/jessfraz/netscan/workflows/make%20all/badge.svg)](https://github.com/jessfraz/netscan/actions?query=workflow%3A%22make+all%22)
[![make-image](https://github.com/jessfraz/netscan/workflows/make%20image/badge.svg)](https://github.com/jessfraz/netscan/actions?query=workflow%3A%22make+image%22)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://godoc.org/github.com/jessfraz/netscan)

Scan a network for ports that are open on an ip/ip range, and
ips that are in use on that network.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [Installation](#installation)
    - [Binaries](#binaries)
    - [Via Go](#via-go)
- [Usage](#usage)
    - [Examples](#examples)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->


## Installation

#### Binaries

For installation instructions from binaries please visit the [Releases Page](https://github.com/jessfraz/netscan/releases).

#### Via Go

```console
$ go get github.com/jessfraz/netscan
```

## Usage

```console
netscan -  Scan network ips and ports.

Usage: netscan <command>

Flags:

  -d, --debug    enable debug logging (default: false)
  -p, --ports    Ports to scan (ex. 80-443 or 80,443,8080 or 1-20,22,80-443) (default "80,443,8001,9001") 
  --proto        protocol to use (can be set more than once) (default "tcp")
  -t, --timeout  timeout for ping of port (default: 1s)

Commands:

  version  Show the version information.
```

#### Examples

```console
# for a cidr
$ netscan 192.168.0.1/24

# for a single ip
$ netscan 192.168.104.30
```

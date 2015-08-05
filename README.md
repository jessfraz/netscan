# netscan

Scan a network for ports that are open on an ip/ip range, and
ips that are in use on that network.

```console
$ netscan -h
NAME:
   netscan - Scan network ips and ports.

USAGE:
   netscan [global options] command [command options] [arguments...]
   
VERSION:
   v0.1.0
   
AUTHOR(S):
   @jfrazelle <no-reply@butts.com> 
   
COMMANDS:
   help, h  Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --debug, -d      run in debug mode
   --timeout, -t "1s"   override timeout used for check
   --port, -p "1-1000"  port range to check
   --proto "tcp,udp"    protocol/s to check
   --help, -h       show help
   --version, -v    print the version
```

Examples:

```console
# for a cidr
$ netscan 192.168.0.1/24

# for a single ip
$ netscan 192.168.104.30
```

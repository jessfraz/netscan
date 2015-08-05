package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

const (
	VERSION = "v0.1.0"
)

var (
	protos = []string{"tcp", "udp"}
)

// preload initializes any global options and configuration
// before the main or sub commands are run
func preload(context *cli.Context) error {
	if context.GlobalBool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}
	return nil
}

func isReachable(proto, addr string, timeout time.Duration) bool {
	c, err := net.DialTimeout(proto, addr, timeout)
	if err == nil {
		c.Close()
		return true
	}
	return false
}

func scanIP(ip net.IP, beginPort, endPort int, timeout time.Duration) {
	logrus.Infof("scanning ip: %v", ip)

	for _, proto := range protos {
		for port := beginPort; port <= endPort; port++ {
			addr := fmt.Sprintf("%s:%d", ip.String(), port)
			logrus.Debugf("scanning addr: %s://%s", proto, addr)
			if isReachable(proto, addr, timeout) {
				logrus.Infof("%s://%s is alive and reachable", proto, addr)
			}
		}
	}
}

func scan(s string, beginPort, endPort int, timeout time.Duration) error {
	ip, ipNet, err := net.ParseCIDR(s)
	if err != nil {
		ip = net.ParseIP(s)
		scanIP(ip, beginPort, endPort, timeout)
		return nil
	}
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); incIP(ip) {
		scanIP(ip, beginPort, endPort, timeout)
	}
	return nil
}

func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func parsePortRange(ports string) (begin, end int, err error) {
	p := strings.SplitN(ports, "-", 2)
	if len(p) < 2 {
		logrus.Debugf("Looks like only one port %q was given for ports.", ports)
		begin, err = strconv.Atoi(p[0])
		end = begin
		return begin, end, err
	}

	begin, err = strconv.Atoi(p[0])
	if err != nil {
		return begin, end, err
	}
	end, err = strconv.Atoi(p[1])
	if err != nil {
		return begin, end, err
	}

	if begin > end {
		return begin, end, fmt.Errorf("End port can not be greater than the beginning port: %d > %d", end, begin)
	}

	return begin, end, err
}

func main() {
	app := cli.NewApp()
	app.Name = "netscan"
	app.Version = VERSION
	app.Author = "@jfrazelle"
	app.Email = "no-reply@butts.com"
	app.Usage = "Scan network ips and ports."
	app.Before = preload
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "run in debug mode",
		},
		cli.StringFlag{
			Name:  "timeout, t",
			Value: "1s",
			Usage: "override timeout used for check",
		},
		cli.StringFlag{
			Name:  "port, p",
			Value: "0-1000",
			Usage: "port range to check",
		},
	}
	app.Action = func(c *cli.Context) {
		if len(c.Args()) == 0 {
			logrus.Errorf("Pass an ip or cidr, ex: 192.168.104.1/24")
			cli.ShowAppHelp(c)
			return
		}

		timeout, err := time.ParseDuration(c.String("timeout"))
		if err != nil {
			logrus.Error(err)
			return
		}

		beginPort, endPort, err := parsePortRange(c.String("port"))
		if err != nil {
			logrus.Error(err)
			return
		}

		err = scan(c.Args().First(), beginPort, endPort, timeout)
		if err != nil {
			logrus.Error(err)
			return
		}
	}
	app.Run(os.Args)
}

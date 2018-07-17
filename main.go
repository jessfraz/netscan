package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/genuinetools/pkg/cli"
	"github.com/jessfraz/netscan/version"
	"github.com/sirupsen/logrus"
)

var (
	timeout   time.Duration
	ports     intSlice
	protocols stringSlice

	defaultProtocols = stringSlice{"tcp"}
	defaultPorts     = intSlice{80, 443, 8001, 9001}
	originalPorts    string

	debug bool
)

// stringSlice is a slice of strings
type stringSlice []string

// implement the flag interface for stringSlice
func (s *stringSlice) String() string {
	return fmt.Sprintf("%s", *s)
}
func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// intSlice is a slice of ints
type intSlice []int

// implement the flag interface for intSlice
func (i *intSlice) String() (out string) {
	for k, v := range *i {
		if k < len(*i)-1 {
			out += fmt.Sprintf("%d,", v)
		} else {
			out += fmt.Sprintf("%d", v)
		}
	}
	return out
}

func (i *intSlice) Set(value string) error {
	originalPorts = value

	// Set the default if nothing was given.
	if len(value) <= 0 {
		*i = defaultPorts
		return nil
	}

	// Split on "," for individual ports and ranges.
	r := strings.Split(value, ",")
	for _, pr := range r {
		// Split on "-" to denote a range.
		if strings.Contains(pr, "-") {
			p := strings.SplitN(pr, "-", 2)
			begin, err := strconv.Atoi(p[0])
			if err != nil {
				return err
			}
			end, err := strconv.Atoi(p[1])
			if err != nil {
				return err
			}
			if begin > end {
				return fmt.Errorf("End port can not be greater than the beginning port: %d > %d", end, begin)
			}
			for port := begin; port <= end; port++ {
				*i = append(*i, port)
			}

			continue
		}

		// It is not a range just parse the port
		port, err := strconv.Atoi(pr)
		if err != nil {
			return err
		}
		*i = append(*i, port)
	}

	return nil
}

func main() {
	// Create a new cli program.
	p := cli.NewProgram()
	p.Name = "netscan"
	p.Description = "Scan network ips and ports"

	// Set the GitCommit and Version.
	p.GitCommit = version.GITCOMMIT
	p.Version = version.VERSION

	// Setup the global flags.
	p.FlagSet = flag.NewFlagSet("global", flag.ExitOnError)
	p.FlagSet.DurationVar(&timeout, "timeout", time.Second, "timeout for ping of port")
	p.FlagSet.DurationVar(&timeout, "t", time.Second, "timeout for ping of port")

	p.FlagSet.Var(&ports, "ports", fmt.Sprintf("Ports to scan (ex. 80-443 or 80,443,8080 or 1-20,22,80-443) (default %q)", defaultPorts.String()))
	p.FlagSet.Var(&ports, "p", fmt.Sprintf("Ports to scan (ex. 80-443 or 80,443,8080 or 1-20,22,80-443) (default %q)", defaultPorts.String()))

	p.FlagSet.Var(&protocols, "proto", `protocol to use (can be set more than once) (default "tcp")`)

	p.FlagSet.BoolVar(&debug, "d", false, "enable debug logging")
	p.FlagSet.BoolVar(&debug, "debug", false, "enable debug logging")

	// Set the before function.
	p.Before = func(ctx context.Context) error {
		// Set the log level.
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}

		if p.FlagSet.NArg() < 1 {
			return errors.New("Pass an ip or cidr, ex: 192.168.104.1/24")
		}

		// Set the default ports.
		if len(ports) < 1 {
			ports = defaultPorts
		}

		if len(protocols) < 1 {
			protocols = defaultProtocols
		}

		return nil
	}

	// Set the main program action.
	p.Action = func(ctx context.Context, args []string) error {
		// On ^C, or SIGTERM handle exit.
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		signal.Notify(c, syscall.SIGTERM)
		go func() {
			for sig := range c {
				logrus.Infof("Received %s, exiting.", sig.String())
				os.Exit(0)
			}
		}()

		logrus.Infof("Scanning on %s using protocols (%s) over ports %s", args[0], strings.Join(protocols, ","), ports.String())

		var (
			ip    net.IP
			ipnet *net.IPNet
			err   error
		)
		if !strings.Contains(args[0], "/") {
			// We got an ip not a CIDR range.
			ip = net.ParseIP(args[0])
			scanIP(ip)
			return nil
		}

		ip, ipnet, err = net.ParseCIDR(args[0])
		if err != nil {
			return err
		}

		for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
			scanIP(ip)
		}

		return nil
	}

	// Run our program.
	p.Run()
}

func scanIP(ip net.IP) {
	var wg sync.WaitGroup
	for _, port := range ports {
		for _, proto := range protocols {
			addr := fmt.Sprintf("%s:%d", ip, port)
			logrus.Debugf("Scannng %s://%s", proto, addr)

			wg.Add(1)
			go func(proto, addr string) {
				defer wg.Done()

				c, err := net.DialTimeout(proto, addr, timeout)
				if err == nil {
					c.Close()
					logrus.Infof("%s://%s is alive and reachable", proto, addr)
				}
			}(proto, addr)
		}
	}
	wg.Wait()
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

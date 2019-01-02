package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/genuinetools/pkg/cli"
	"github.com/jessfraz/netscan/pkg/scanner"
	"github.com/jessfraz/netscan/version"
	"github.com/sirupsen/logrus"
)

var (
	timeout         time.Duration
	ports           intSlice
	protocols       stringSlice
	parallelRunners int

	defaultProtocols       = stringSlice{"tcp"}
	defaultPorts           = intSlice{80, 443, 8001, 9001}
	defaultParallelRunners = 100
	originalPorts          string

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
				return fmt.Errorf("end port can not be greater than the beginning port: %d > %d", end, begin)
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
	log := logrus.New()

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

	p.FlagSet.IntVar(&parallelRunners, "r", defaultParallelRunners, fmt.Sprintf("Maximum amount of parallel runners (default %q)", defaultParallelRunners))
	p.FlagSet.IntVar(&parallelRunners, "runners", defaultParallelRunners, fmt.Sprintf("Maximum amount of parallel runners (default %q)", defaultParallelRunners))

	p.FlagSet.Var(&protocols, "proto", `protocol to use (can be set more than once) (default "tcp")`)

	p.FlagSet.BoolVar(&debug, "d", false, "enable debug logging")
	p.FlagSet.BoolVar(&debug, "debug", false, "enable debug logging")

	// Set the before function.
	p.Before = func(ctx context.Context) error {
		// Set the log level.
		if debug {
			log.SetLevel(logrus.DebugLevel)
		}

		if p.FlagSet.NArg() < 1 {
			return errors.New("pass an ip or cidr, ex: 192.168.104.1/24")
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
				log.Infof("Received %s, exiting.", sig.String())
				os.Exit(0)
			}
		}()

		log.Infof("Scanning on %s using protocols (%s) over ports %s", args[0], strings.Join(protocols, ","), ports.String())

		scan := scanner.NewScanner(scanner.WithTimeout(timeout), scanner.WithProtocols(protocols), scanner.WithParallelRunners(parallelRunners))

		var err error
		if !strings.Contains(args[0], "/") {
			err = scan.AddIP(args[0])
		} else {
			err = scan.AddCIDR(args[0])
		}
		if err != nil {
			return err
		}

		scan.SetPorts(ports)
		scan.ScanToLogger(log)

		return nil
	}

	// Run our program.
	p.Run()
}

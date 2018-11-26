package scanner

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// Scanner allows to port scan multiple IPs for multiple ports
type Scanner struct {
	ips       []net.IP
	timeout   time.Duration
	ports     []int
	protocols []string
}

// AddressSet is a set of the IP port and protocol of a responding service
type AddressSet struct {
	IP       net.IP
	Port     int
	Protocol string
}

// NewScanner gives a new Scanner instance
func NewScanner(options ...func(*Scanner)) *Scanner {
	scanner := Scanner{
		timeout:   time.Second,
		ports:     []int{80, 443, 8001, 9001},
		protocols: []string{"tcp"},
	}

	for _, option := range options {
		option(&scanner)
	}

	return &scanner
}

// WithTimeout is used as an option in NewScanner to set the timeout for the port dial
func WithTimeout(timeout time.Duration) func(*Scanner) {
	return func(s *Scanner) {
		s.timeout = timeout
	}
}

// WithProtocols is used as an option in NewScanner to set the protocols to test the ports with
func WithProtocols(protocols []string) func(*Scanner) {
	return func(s *Scanner) {
		s.protocols = protocols
	}
}

// WithPorts is used as an option in NewScanner to set the ports to scan
func WithPorts(ports []int) func(*Scanner) {
	return func(s *Scanner) {
		s.SetPorts(ports)
	}
}

// AddCIDR adds all IPs in a CIDR notation to the list of IPs to be scanned
func (s *Scanner) AddCIDR(cidr string) error {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}

	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		s.ips = append(s.ips, copyIP(ip))
	}

	return nil
}

// AddIP adds a single IP to the list of IPs to be scanned
func (s *Scanner) AddIP(ip string) error {
	netIP := net.ParseIP(ip)
	if netIP == nil {
		return fmt.Errorf("%s is not a valid IP", ip)
	}

	s.ips = append(s.ips, netIP)
	return nil
}

// SetPorts sets the port to be scanned
func (s *Scanner) SetPorts(ports []int) {
	s.ports = ports
}

// Scan performs a network scan and returns the responding addressses
func (s *Scanner) Scan() []AddressSet {
	results := []AddressSet{}
	resultsMutex := sync.Mutex{}
	var wg sync.WaitGroup
	for _, ip := range s.ips {
		for _, port := range s.ports {
			for _, proto := range s.protocols {
				addr := fmt.Sprintf("%s:%d", ip, port)

				wg.Add(1)
				go func(proto, addr string) {
					defer wg.Done()

					c, err := net.DialTimeout(proto, addr, s.timeout)
					if err == nil {
						c.Close()
						resultsMutex.Lock()
						results = append(results, AddressSet{
							IP:       copyIP(ip),
							Port:     port,
							Protocol: proto,
						})
						resultsMutex.Unlock()
					}
				}(proto, addr)
			}
		}
	}
	wg.Wait()

	return results
}

type logger interface {
	Debugf(format string, a ...interface{})
	Infof(format string, a ...interface{})
}

// ScanToLogger perforn a scan and outputs the results to a logger interface
func (s *Scanner) ScanToLogger(log logger) {
	var wg sync.WaitGroup
	for _, ip := range s.ips {
		for _, port := range s.ports {
			for _, proto := range s.protocols {
				addr := fmt.Sprintf("%s:%d", ip, port)
				log.Debugf("Scannng %s://%s", proto, addr)

				wg.Add(1)
				go func(proto, addr string) {
					defer wg.Done()

					c, err := net.DialTimeout(proto, addr, s.timeout)
					if err == nil {
						c.Close()
						log.Infof("%s://%s is alive and reachable", proto, addr)
					}
				}(proto, addr)
			}
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

func copyIP(ip net.IP) net.IP {
	dup := make(net.IP, len(ip))
	copy(dup, ip)
	return dup
}

package runner

import (
	"net"

	"github.com/ExploitSuite/cdncheck"
)

type CDNDetector struct {
	client *cdncheck.Client
}

func NewCDNDetector() *CDNDetector {
	return &CDNDetector{
		client: cdncheck.New(),
	}
}

func (d *CDNDetector) Detect(input string) (bool, string, string, error) {
	ip := net.ParseIP(input)
	if ip != nil {
		return d.client.Check(ip)
	}
	return d.client.CheckDomainWithFallback(input)
}

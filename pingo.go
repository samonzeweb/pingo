package pingo

import "net"

const (
	IP   = "ip"
	IPv4 = "ip4"
	IPv6 = "ip6"
)

const (
	protocolICMPv4 = 1
	protocolICMPv6 = 58
)

func resolve(target, ipVersion string) (*net.IPAddr, error) {
	var address *net.IPAddr

	// IP address was given
	ip := net.ParseIP(target)
	if ip != nil {
		return &net.IPAddr{IP: ip}, nil
	}

	// Not a valid IP, perhaps a host
	address, err := net.ResolveIPAddr(ipVersion, target)
	if err != nil {
		return nil, err
	}

	return address, nil
}

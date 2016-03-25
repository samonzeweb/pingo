// Package pingo is a simple go package to do ICMP ping (echo) requests
package pingo

import (
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

const (
	// IP resolve with IPv4 or IPv6
	IP = "ip"
	// IPv4 only
	IPv4 = "ip4"
	// IPv6 only
	IPv6 = "ip6"
)

var (
	// ErrTimeOut : no reply before specified timeout
	ErrTimeOut = errors.New("Ping : time out")
	// ErrDestinationUnreachable : self explanatory
	ErrDestinationUnreachable = errors.New("Ping : destination unreachable")
	// ErrTimeExceeded : often a TTL, fragmentation, ... trouble
	ErrTimeExceeded = errors.New("Ping : time exceeded")
)

const (
	protocolICMPv4 = 1
	protocolICMPv6 = 58
)

const (
	networkIPICMPv4 = "ip4:icmp"
	networkIPICMPv6 = "ip6:ipv6-icmp"
)

var payload []byte

func init() {
	payload = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

// SimplePing make an icmp echo request and wait for an ICMP echo reply (before timeout)
//
// The target is either an hostname, or an IPv4 ou IPv6 address.
// Use the constants IP, IPv4 or IPv6 for ipVersion to specify wich IP version will be
// used when resolving hostname (useless with an IP address).
//
// If the remote host replies before timeout occurs, SimplePing return the trip duration.
// In any other case the duration is zero and an error is returned. For a simple timeout
// the error is ErrTimeOut.
func SimplePing(target string, ipVersion string, timeout time.Duration) (time.Duration, error) {
	var err error

	// Address resolution (if needed)
	address, err := resolve(target, ipVersion)
	if err != nil {
		return 0, err
	}

	// IPv4 or IPv6 ?
	var protocol int
	var network string
	var localAddress string
	var echoMessage icmp.Type
	if address.IP.To4() != nil {
		protocol = protocolICMPv4
		network = networkIPICMPv4
		echoMessage = ipv4.ICMPTypeEcho
		localAddress = "0.0.0.0"
	} else {
		protocol = protocolICMPv6
		network = networkIPICMPv6
		echoMessage = ipv6.ICMPTypeEchoRequest
		localAddress = "::"
	}

	// Create ICMP Connection
	conn, err := icmp.ListenPacket(network, localAddress)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	// Build echo request
	echo := &icmp.Echo{ID: os.Getpid(), Seq: 0, Data: payload}
	message := &icmp.Message{Type: echoMessage,
		Code: 0,
		Body: echo}
	fullMessage, err := message.Marshal(nil)
	if err != nil {
		return 0, err
	}

	// Send echo request
	err = conn.SetWriteDeadline(time.Now().Add(timeout))
	if err != nil {
		return 0, err
	}

	sendTime := time.Now()
	_, err = conn.WriteTo(fullMessage, address)
	if err != nil {
		return 0, readWriteError(err)
	}
	remainingDuration := timeout - time.Now().Sub(sendTime)
	if remainingDuration <= 0 {
		return 0, ErrTimeOut
	}

	// Wait and read reply
	err = conn.SetReadDeadline(time.Now().Add(remainingDuration))
	if err != nil {
		return 0, err
	}
	buffer := make([]byte, 512)
	_, _, err = conn.ReadFrom(buffer)
	if err != nil {
		return 0, readWriteError(err)
	}
	totalDuration := time.Since(sendTime)
	if totalDuration > timeout {
		return 0, ErrTimeOut
	}

	// Parse reply
	message, err = icmp.ParseMessage(protocol, buffer)
	if err != nil {
		return 0, err
	}
	switch icmpMessage := message.Body.(type) {
	case *icmp.Echo:
		// Success :)
		return totalDuration, nil
	case *icmp.DstUnreach:
		return 0, ErrDestinationUnreachable
	case *icmp.TimeExceeded:
		return 0, ErrTimeExceeded
	default:
		return 0, fmt.Errorf("Unexpected ICMP message type : %T", icmpMessage)
	}
}

// resolve hostname if needed, or parse the IP string
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

// readWriteError replace error with the package ErrTimeOut one if it is
// a read or write timeout error.
func readWriteError(err error) error {
	netError, ok := err.(net.Error)
	if ok && netError.Timeout() {
		return ErrTimeOut
	}
	return err
}

//go:build freebsd

package tunsetup

import (
	"fmt"
	"net"
	"os/exec"
)

func configureTunAddrImpl(tunName, ipv4, ipv6 string) error {
	// IPv4
	ip, mask := parseCIDR(ipv4)
	if err := exec.Command("ifconfig", tunName, "inet", ip, mask, "up").Run(); err != nil {
		return fmt.Errorf("failed to add IPv4 address to %s: %w", tunName, err)
	}

	// IPv6
	ip6, prefix := parseIPv6Prefix(ipv6)
	if err := exec.Command("ifconfig", tunName, "inet6", ip6, "prefixlen", prefix, "up").Run(); err != nil {
		return fmt.Errorf("failed to add IPv6 address to %s: %w", tunName, err)
	}

	return nil
}

// parseCIDR extracts IP and netmask string
func parseCIDR(cidr string) (string, string) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return cidr, "255.255.255.0"
	}
	mask := net.IP(ipnet.Mask).String()
	return ip.String(), mask
}

func parseIPv6Prefix(cidr string) (string, string) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return cidr, "64"
	}
	prefix, _ := ipnet.Mask.Size()
	return ip.String(), fmt.Sprintf("%d", prefix)
}

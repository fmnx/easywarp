//go:build windows

package tunsetup

import (
	"fmt"
	"os/exec"
)

func configureTunAddrImpl(tunName, ipv4, ipv6 string) error {
	if err := exec.Command("netsh", "interface", "ipv4", "set", "address", tunName, "static", ipv4).Run(); err != nil {
		return fmt.Errorf("failed to add ipv4 address for tun device %s: %w", tunName, err)
	}

	if err := exec.Command("netsh", "interface", "ipv6", "add", "address", tunName, ipv6).Run(); err != nil {
		return fmt.Errorf("failed to add ipv6 address for tun device %s: %w", tunName, err)
	}

	if err := exec.Command("netsh", "interface", "ipv4", "set", "dnsservers", fmt.Sprintf("name=%s", tunName),
		"static", "address=1.1.1.1", "register=none", "validate=no").Run(); err != nil {
		return fmt.Errorf("failed to set dns server for tun %s: %w", tunName, err)
	}

	return nil
}

//go:build linux

package tunsetup

import (
	"fmt"
	"os/exec"
)

func configureTunAddrImpl(tunName, ipv4, ipv6 string) error {
	if err := exec.Command("ip", "link", "set", tunName, "up").Run(); err != nil {
		return fmt.Errorf("failed to set %s up: %w", tunName, err)
	}
	if err := exec.Command("ip", "addr", "add", ipv4, "dev", tunName).Run(); err != nil {
		return fmt.Errorf("failed to add IPv4 address to %s: %w", tunName, err)
	}
	if err := exec.Command("ip", "-6", "addr", "add", ipv6, "dev", tunName).Run(); err != nil {
		return fmt.Errorf("failed to add IPv6 address to %s: %w", tunName, err)
	}
	return nil
}

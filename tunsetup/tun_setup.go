package tunsetup

func ConfigureTunAddr(tunName, ipv4, ipv6 string) error {
	return configureTunAddrImpl(tunName, ipv4, ipv6)
}

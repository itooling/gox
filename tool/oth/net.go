package oth

import (
	"net"
)

func GetIP() ([]string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	ips := make([]string, 0)

	for _, v := range addrs {
		if ipnet, ok := v.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}

		}
	}
	return ips, nil
}

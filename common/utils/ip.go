package utils

import (
	"fmt"
	"net"
)

func GetIPRange(ip_start, ip_end string) ([]string, error) {
	var ips []string
	ip_s, ipnet_s, err := net.ParseCIDR(ip_start)
	if err != nil {
		return nil, err
	}
	ip_e, ipnet_e, err := net.ParseCIDR(ip_end)
	if err != nil {
		return nil, err
	}
	if ipnet_s.Mask.String() != ipnet_e.Mask.String() {
		return nil, fmt.Errorf("%s and %s are not in the same subnet", ip_start, ip_end)
	}
	for ip := ip_s; ipnet_s.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
		if ip.Equal(ip_e) {
			break
		}
	}
	return ips, nil

}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

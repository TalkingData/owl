package utils

import (
	"fmt"
	"net"
)

func GenIpByCidrRange(cidrStart, cidrEnd string) ([]string, error) {
	var ips []string

	// 对于开始地址比结束地址大的，不予处理
	if cidrStart > cidrEnd {
		return nil, fmt.Errorf("%s greater than %s", cidrStart, cidrEnd)
	}

	ipStart, ipNetStart, err := net.ParseCIDR(cidrStart)
	if err != nil {
		return nil, err
	}

	ipEnd, ipNetEnd, err := net.ParseCIDR(cidrEnd)
	if err != nil {
		return nil, err
	}

	// 对于开始地址与结束地址所属网络不一样的，不予处理
	if ipNetStart.Mask.String() != ipNetEnd.Mask.String() {
		return nil, fmt.Errorf("%s and %s are not in the same subnet", cidrStart, cidrEnd)
	}

	// 遍历生成ip地址
	for ip := ipStart; ipNetStart.Contains(ip); innerFunc(ip) {
		ips = append(ips, ip.String())
		if ip.Equal(ipEnd) {
			break
		}
	}

	return ips, nil
}

func innerFunc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

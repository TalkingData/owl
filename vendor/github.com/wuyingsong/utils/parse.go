package utils

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func ParsePort(portString string) []int {
	ports := []int{}
	m := make(map[int]struct{})
	for _, f1 := range strings.Split(portString, ",") {
		if strings.Contains(f1, "-") {
			f2 := strings.Split(f1, "-")
			if len(f2) != 2 {
				continue
			}
			s, err := strconv.Atoi(strings.TrimSpace(f2[0]))
			if err != nil {
				continue
			}
			e, err := strconv.Atoi(strings.TrimSpace(f2[1]))
			if err != nil {
				continue
			}
			if s >= e || s <= 0 {
				continue
			}
			for ; s < e; s++ {
				if _, ok := m[s]; ok {
					continue
				}
				ports = append(ports, s)
				m[s] = struct{}{}
			}
		} else {
			port, err := strconv.Atoi(strings.TrimSpace(f1))
			if err != nil {
				continue
			}
			if _, ok := m[port]; ok {
				continue
			}
			ports = append(ports, port)
			m[port] = struct{}{}
		}
	}
	return ports
}

func ParseCommandArgs(s string) []string {
	rd := bufio.NewReader(strings.NewReader(s))
	fields := []string{}
	var flag byte
	for {
		field, err := rd.ReadString(32)
		if err != nil {
			if err == io.EOF && len(field) > 0 {
				fields = append(fields, strings.TrimSpace(field))
			}
			break
		}
		if field == " " {
			continue
		}

		if strings.Contains(field, string(34)) { //双引号
			flag = 34
		}
		if flag != 0 {
			s, err := rd.ReadString(flag)
			if err != nil {
				break
			}
			field = strings.Trim(fmt.Sprintf("%s %s", field, s), string(34))
			flag = 0
		}
		fields = append(fields, strings.TrimSpace(field))
	}
	return fields
}

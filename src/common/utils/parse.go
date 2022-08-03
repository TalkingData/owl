package utils

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

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

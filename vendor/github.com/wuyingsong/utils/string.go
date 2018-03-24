package utils

import (
	"strings"
)

func TrimSpaceAndNewLine(str string) string {
	str = strings.TrimSpace(str)
	return strings.Trim(str, "\n")
}

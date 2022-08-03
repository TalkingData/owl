package utils

import "fmt"

func OsArch(os, arch string) string {
	if len(os) < 1 || len(arch) < 1 {
		return ""
	}

	return fmt.Sprintf("%s_%s", os, arch)
}

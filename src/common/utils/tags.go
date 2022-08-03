package utils

import (
	"fmt"
	"sort"
	"strings"
)

func Tags2String(tags map[string]string) string {
	if len(tags) == 0 {
		return ""
	}

	keys := make([]string, len(tags))
	i := 0
	for k := range tags {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	ret := ""
	for _, k := range keys {
		ret += fmt.Sprintf("%s=%s,", strings.TrimSpace(k), strings.TrimSpace(tags[k]))
	}
	return strings.Trim(ret, ",")
}

package utils

import "fmt"

func Bytes2Human(num float64) string {
	if num < 1000.00 {
		return fmt.Sprintf("%3.2f", num)
	}
	suffix := []string{"", "K", "M", "G", "T", "P", "E", "Z"}
	for _, unit := range suffix {
		if num < 1000.00 {
			return fmt.Sprintf("%3.2f%s%s", num, unit, "B")
		}
		num /= 1000.00
	}
	return fmt.Sprintf("%.2f%s%s", num, "Y", "B")
}

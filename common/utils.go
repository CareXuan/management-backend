package common

import (
	"fmt"
	"strings"
)

func DecimalToHex(decimal int) string {
	hex := fmt.Sprintf("%x", decimal)
	if len(hex) < 4 {
		hex = strings.Repeat("0", 4-len(hex)) + hex
	}
	return hex
}

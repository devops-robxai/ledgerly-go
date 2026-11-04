package billing

import "fmt"

// FormatUSD formats cents as a US dollar string, e.g. 24900 -> "$249.00".
func FormatUSD(cents int) string {
	neg := cents < 0
	if neg {
		cents = -cents
	}
	dollars := cents / 100
	rem := cents % 100
	s := fmt.Sprintf("$%d.%02d", dollars, rem)
	if neg {
		return "-" + s
	}
	return s
}

package helper

import (
	"strings"
)

func PadRight(s string, n int) string {
	return s + strings.Repeat(" ", n)
}


// Align aligns (duh) a multiline string.
// It splits each line by provided separator sep
// and makes sure the left side of each line is of equal length
func Align(text string, sep string) string {
	lines := strings.Split(text, "\n")
	longest := 0
	aligned := make([]string, len(lines))
	for _, l := range lines {
		parts := strings.Split(l, "-")
		if len(parts[0]) > longest {
			longest = len(parts[0])
		}
	}

	for i, l := range lines {
		parts := strings.Split(l, "-")
		aligned[i] = PadRight(parts[0], longest-len(parts[0])) + sep + parts[1]
	}

	return strings.Join(aligned, "\n")
}

package helper

import (
	"fmt"
	"strings"
)

func PadRight(s string, n int) string {
	return s + strings.Repeat(" ", n)
}

func AddModifier(s string, mod string) string {
	return fmt.Sprintf("[%v]%v[-:-:-]", mod, s)
}

// Columns splits each line with provided separator sep
// and makes sure the left side of each line is of equal length.
// Then the modifiers modL and modR are applied to the left and right columns
func Columns(text string, sep string, modL string, modR string) string {
	lines := strings.Split(text, "\n")
	longest := 0
	aligned := make([]string, len(lines))
	for _, l := range lines {
		parts := strings.Split(l, sep)
		if len(parts[0]) > longest {
			longest = len(parts[0])
		}
	}

	for i, l := range lines {
		parts := strings.Split(l, sep)
		aligned[i] = PadRight(parts[0], longest-len(parts[0]))
		if modL != "" {
			aligned[i] = AddModifier(aligned[i], modL)
		}
		if len(parts) > 1 {
			if modR != "" {
				parts[1] = AddModifier(parts[1], modR)
			}
			aligned[i] += sep + parts[1]
		}
	}

	return strings.Join(aligned, "\n")
}

// PadLinesRight pads all the lines of string s on the right
// So that all strings are the same length
func PadLinesRight(s string) string {
	lines := strings.Split(s, "\n")
	longest := len(lines[0])
	for _, l := range lines {
		if len(l) > longest {
			longest = len(l)
		}
	}

	for i, l := range lines {
		lines[i] = PadRight(l, longest-len(l))
	}

	return strings.Join(lines, "\n")
}

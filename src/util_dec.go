package spejt

import (
	term "github.com/tj/go/term"
)

func DashBorder() string {
	width, _ := term.Size()
	var toPrint string
	for n := 1; n <= width; n++ {
		toPrint = toPrint + "-"
	}
	return toPrint
}

func DashBorder2(text string) string {
	width, _ := term.Size()
	var toPrint = "--- " + text + " -"
	var toPrintLen = len([]rune(toPrint))
	for n := 1; n <= width-toPrintLen; n++ {
		toPrint = toPrint + "-"
	}
	toPrint = toPrint + "\n"
	return toPrint
}

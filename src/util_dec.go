package shko

import (
	"fmt"
)

func DashBorder(symb string) string {
	var toPrint string
	for n := 1; n <= termWidth; n++ {
		toPrint = toPrint + symb
	}
	toPrint = toPrint + "\n"
	return toPrint
}

func DashBorder2(text, symb string, before int) string {
	var toPrint = ""
	for n := 1; n <= before; n++ {
		toPrint += symb
	}
	if text != "" {
		toPrint += " " + text + " " + symb
	}
	var toPrintLen = len([]rune(toPrint))
	for n := 1; n <= termWidth-toPrintLen; n++ {
		toPrint = toPrint + symb
	}
	toPrint = toPrint + "\n"
	return toPrint
}

func Print(stl Style, fg Color, bg Color, toPrint string) {
	SetStyle(stl, fg, bg)
	fmt.Print(toPrint)
	ResetStyle()
}

func PrintWait(toPrint string) {
	SetStyle(HighLight, Black, White)
	fmt.Print(toPrint)
	ResetStyle()
	GetChar()
}

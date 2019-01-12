package shko

import (
	"fmt"

	colorz "github.com/bresilla/shko/colorz"
	keyboard "github.com/bresilla/shko/keyboard"
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

func Print(stl colorz.Style, fg colorz.Color, bg colorz.Color, toPrint string) {
	colorz.SetStyle(stl, fg, bg)
	fmt.Print(toPrint)
	colorz.ResetStyle()
}

func PrintWait(toPrint string) {
	colorz.SetStyle(colorz.HighLight, colorz.Black, colorz.White)
	fmt.Print(toPrint)
	colorz.ResetStyle()
	keyboard.GetChar()
}

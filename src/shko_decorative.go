package shko

import (
	"fmt"

	keyboard "github.com/bresilla/shko/keyboard"
	t "github.com/bresilla/shko/term"
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

func Print(stl t.Style, fg t.Color, bg t.Color, toPrint string) {
	t.SetStyle(stl, fg, bg)
	fmt.Print(toPrint)
	t.ResetStyle()
}

func PrintWait(toPrint string) {
	t.SetStyle(t.HighLight, t.Black, t.White)
	fmt.Print(toPrint)
	t.ResetStyle()
	keyboard.GetChar()
}

package shko

import (
	"fmt"
	"log"

	t "github.com/bresilla/shko/term"
	"github.com/peterh/liner"
	"github.com/tj/go/term"
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

func StatusWrite(toWrite string) {
	term.MoveTo(0, termHeight+1)
	cleanLine(0)
	term.MoveTo(0, termHeight+1)
	fmt.Print(toWrite)
	fmt.Print(" ")
}

func StatusRead(toWrite, defaultStr string) (text string) {
	line := liner.NewLiner()
	defer line.Close()
	line.SetCtrlCAborts(true)
	term.MoveTo(0, termHeight+1)
	if name, err := line.PromptWithSuggestion(toWrite+": ", defaultStr, -1); err == nil {
		text = name
	} else if err == liner.ErrPromptAborted {
		log.Print("Aborted")
	} else {
		log.Print("Error reading line: ", err)
	}
	return
}

func cleanLine(minus int) {
	for i := 0; i < termWidth-minus; i++ {
		fmt.Print(" ")
	}
}

func Print(stl t.Style, fg t.Color, bg t.Color, toPrint string) {
	t.SetStyle(stl, fg, bg)
	fmt.Print(toPrint)
	t.ResetStyle()
}

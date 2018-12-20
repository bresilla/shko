package spejt

import (
	"fmt"

	term "github.com/buger/goterm"
)

func SelectInList(selected int, file []File) {
	term.MoveCursor(0, 0)
	term.Flush()
	term.Clear()
	for num, el := range file {
		if num < selected {
			colorList(el, false)
		} else if num == selected {
			colorList(el, true)
		} else if num > selected {
			colorList(el, false)
		}
	}
	if len(file) == 0 {
		Print(HighLight, Black, Cyan, "\t nothing to show ")
	}
}

func colorList(file File, active bool) {
	var name string
	if file.IsDir {
		name = "»  " + file.Name
	} else {
		name = "♦  " + file.Name
	}
	if file.IsDir && active {
		SetStyle(HighLight, Black, Cyan)
		fmt.Println("\t " + name + " / ")
		ResetStyle()
	} else if file.IsDir && !active {
		SetStyle(HighLight, White, None)
		fmt.Println("\t" + name + "/")
		ResetStyle()
	} else if !file.IsDir && active {
		SetStyle(HighLight, Black, Cyan)
		fmt.Println("\t " + name + " ")
		ResetStyle()
	} else if !file.IsDir && !active {
		SetStyle(HighLight, Grey, None)
		fmt.Println("\t" + name)
		ResetStyle()

	}
}

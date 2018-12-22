package spejt

import (
	"fmt"

	term "github.com/buger/goterm"
)

func find(list []File, actual File) (number, scroll int) {
	for i, el := range list {
		if el.Name == actual.Name {
			if i < term.Height()-1 {
				number = i
				scroll = 0
				break
			} else {
				number = 0
				scroll = i
			}
		} else {
			number = 0
			scroll = 0
		}
	}
	return
}

func SelectInList(selected int, file []File) {
	term.MoveCursor(0, 0)
	term.Flush()
	term.Clear()
	for i, el := range file {
		if i == selected {
			colorList(el, true)
		} else {
			colorList(el, false)
		}
	}
	if len(file) == 0 {
		Print(HighLight, Black, Cyan, "\t nothing to show ")
	}
}

func colorList(file File, active bool) {
	if file.IsDir && active {
		SetStyle(HighLight, Black, Cyan)
		fmt.Println("\t ", file.Icon, "  ", file.Name, " / \t")
		ResetStyle()
	} else if file.IsDir && !active {
		SetStyle(HighLight, White, None)
		fmt.Println("\t", file.Icon, "  ", file.Name, "/ \t")
		ResetStyle()
	} else if !file.IsDir && active {
		SetStyle(HighLight, Black, Cyan)
		fmt.Println("\t ", file.Icon, "  ", file.Name, "\t")
		ResetStyle()
	} else if !file.IsDir && !active {
		SetStyle(Default, Grey, None)
		fmt.Println("\t", file.Icon, "  ", file.Name, "\t")
		ResetStyle()
	}
}

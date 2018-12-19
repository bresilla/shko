package spejt

import (
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
		Print(HighLight, Black, Red, "\t nothing to show ")
	}
}

func colorList(file File, active bool) {
	if file.IsDir && active {
		Print(HighLight, Black, Red, "\t »  "+file.Name+" / ")
	} else if file.IsDir && !active {
		Print(HighLight, White, None, "\t»  "+file.Name+"/")
	} else if !file.IsDir && active {
		Print(HighLight, Black, Red, "\t ♦  "+file.Name+" ")
	} else if !file.IsDir && !active {
		Print(Default, Grey, None, "\t♦  "+file.Name)
	}
}

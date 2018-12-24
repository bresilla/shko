package spejt

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	term "github.com/buger/goterm"
)

func SelectInList(selected int, file []File) {
	term.MoveCursor(0, 0)
	term.Flush()
	term.Clear()
	for i, el := range file {
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 0, 2, '0', tabwriter.AlignRight|tabwriter.Debug)
		fmt.Print("\t")
		if i == selected {
			colorList(el, true)
		} else {
			colorList(el, false)
		}
		fmt.Println()
		ResetStyle()
		w.Flush()
	}
	if len(file) == 0 {
		Print(HighLight, Black, White, "\t nothing to show ")
	}
}

func colorList(file File, active bool) {
	icon := drawIcon(active, showIcons, file)
	name := drawName(active, file)
	chld := drawChildren(showChildren, file)
	all := icon + name + chld
	if file.IsDir {
		Invert(active, HighLight, White, all)
	} else {
		Invert(active, Default, Cyan, all)
	}
}
func drawName(active bool, file File) (back string) {
	spacer := ""
	if active {
		spacer = " "
	}
	if file.IsDir {
		back = file.Name + spacer + "/ \t"
	} else {
		back = file.Name + spacer + "\t"
	}
	return
}
func drawIcon(active, yesno bool, file File) (back string) {
	spacer := ""
	if active {
		spacer = " "
	}
	if yesno {
		back = " " + spacer + file.Other.Icon + "  "
	} else {
		if file.IsDir {
			back = spacer + " >  "
		} else {
			back = spacer + " -  "
		}
	}
	return
}

func drawChildren(yesno bool, file File) (back string) {
	if yesno {
		back = strconv.Itoa(file.Children) + "\t"
	} else {
		back = ""
	}
	return
}

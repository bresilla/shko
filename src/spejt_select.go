package spejt

import (
	"fmt"
	"strconv"
	"time"

	term "github.com/tj/go/term"
)

var (
	tab = 20
)

func SelectInList(selected int, file []File) {
	term.MoveTo(0, 0)
	term.ClearAll()
	if len(file) == 0 {
		fmt.Print("  ")
		Print(HighLight, Black, White, "  nothing to show  ")
		//term.MoveTo(0, termHeight)
		//Print(HighLight, Red, None, DashBorder2("nothing to show"))
	} else {
		for i, el := range file {
			fmt.Print("  ")
			if i == selected {
				colorList(el, true, i)
			} else {
				colorList(el, false, i)
			}
			fmt.Print("\n")
			ResetStyle()
		}
		//term.MoveTo(0, termHeight)
		//Print(HighLight, Red, None, DashBorder2(file[selected].Path))
	}
}

func colorList(file File, active bool, i int) {
	if file.IsDir {
		Invert(active, HighLight, White)
	} else {
		Invert(active, Default, Cyan)
	}
	term.ClearLineEnd()
	tab = drawIcon(active, showIcons, file)
	tab = drawName(active, file)
	tab = drawChildren(showChildren, file, i)
	tab = drawSize(showSize, file, i)
	tab = drawDate(showDate, file, i)
	SetStyle(Default, White, Black)
	term.ClearLineEnd()
}

func drawIcon(active, yesno bool, file File) (tabTurn int) {
	spacer := ""
	if active {
		spacer = " "
	}
	if yesno {
		fmt.Print(" " + spacer + file.Other.Icon + "  ")
	} else {
		if file.IsDir {
			fmt.Print(spacer + " >  ")
		} else {
			fmt.Print(spacer + " -  ")
		}
	}
	tabTurn = 10
	return
}

func drawName(active bool, file File) (tabTurn int) {
	spacer := ""
	if active {
		spacer = " "
	}
	if file.IsDir {
		fmt.Print(file.Name + spacer + "/ ")
	} else {
		fmt.Print(file.Name + spacer + " ")
	}
	tabTurn = tab + 15
	return
}

func drawChildren(yesno bool, file File, i int) (tabTurn int) {
	if yesno {
		term.MoveTo(tab, i+1)
		fmt.Print("\t  " + strconv.Itoa(file.Children) + " ")
		tabTurn = tab + 6
	} else {
		tabTurn = tab
	}
	return tabTurn
}

func drawSize(yesno bool, file File, i int) (tabTurn int) {
	if yesno {
		term.MoveTo(tab, i+1)
		fmt.Print("\t" + file.Other.HumanSize + " ")
		tabTurn = tab + 15
	} else {
		tabTurn = tab
	}
	return tabTurn
}

func drawDate(yesno bool, file File, i int) (tabTurn int) {
	if yesno {
		term.MoveTo(tab, i+1)
		fmt.Print("\t" + file.ModTime.Format(time.RFC822) + " ")
		tabTurn = tab + 25
	} else {
		tabTurn = tab
	}
	return tabTurn
}

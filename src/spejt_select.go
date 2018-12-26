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

func SelectInList(selected int, file, children []File, parent File) {
	term.MoveTo(0, 0)
	term.ClearAll()
	termWidth, termHeight := term.Size()
	topBarSpace := 0
	if topBar {
		topBarSpace = 1
		Print(HighLight, Black, Cyan, DashBorder2(parent.Path, termWidth/2-(len([]rune(parent.Path)))/2))
	}
	if len(file) == 0 {
		fmt.Print("  ")
		Print(HighLight, Black, White, "  nothing to show  ")
	} else {
		var maxSize int64
		for _, el := range children {
			maxSize += el.Size
		}
		for i, el := range file {
			fmt.Print("  ")
			if i == selected {
				colorList(el, true, i+topBarSpace, maxSize)
			} else {
				colorList(el, false, i+topBarSpace, maxSize)
			}
			fmt.Print("\n")
			ResetStyle()
		}
	}
	if statBar {
		term.MoveTo(0, termHeight-1)
		Print(HighLight, Black, Cyan, DashBorder2(parent.Path, termWidth/2-(len([]rune(parent.Path)))/2))
	}
}

func colorList(file File, active bool, i int, maxSize int64) {
	if file.IsDir {
		Invert(active, HighLight, White)
	} else {
		Invert(active, Default, Cyan)
	}
	term.ClearLineEnd()
	tab = drawIcon(active, showIcons, file)
	tab = drawName(active, file)
	tab = drawChildren(showChildren, file, i)
	tab = drawMode(showMode, file, i)
	tab = drawDU(duMode, file, i, maxSize)
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
		fmt.Print("  " + strconv.Itoa(file.Children) + " ")
		tabTurn = tab + 8
	} else {
		tabTurn = tab
	}
	return tabTurn
}

func drawSize(yesno bool, file File, i int) (tabTurn int) {
	if yesno {
		term.MoveTo(tab, i+1)
		fmt.Print(file.Other.HumanSize + " ")
		tabTurn = tab + 12
	} else {
		tabTurn = tab
	}
	return tabTurn
}

func drawDate(yesno bool, file File, i int) (tabTurn int) {
	if yesno {
		term.MoveTo(tab, i+1)
		fmt.Print(file.ModTime.Format(time.RFC822) + " ")
		tabTurn = tab + 25
	} else {
		tabTurn = tab
	}
	return tabTurn
}

func drawMode(yesno bool, file File, i int) (tabTurn int) {
	if yesno {
		term.MoveTo(tab, i+1)
		fmt.Print(file.Mode)
		fmt.Print(" ")
		tabTurn = tab + 12
	} else {
		tabTurn = tab
	}
	return tabTurn
}

func sizeBar(maxSize, size int64) (toPrint string) {
	var (
		load  string
		uload string
	)
	percentage := int(size) * 100 / int(maxSize)
	for i := 1; i <= percentage; i = i + 10 {
		load += "█"
	}
	for i := 1; i <= 10-len([]rune(load)); i++ {
		uload += "░"
	}
	toPrint = "│" + load + uload + "│"
	//toPrint = strconv.Itoa(percentage)
	return
}

func drawDU(yesno bool, file File, i int, maxSize int64) (tabTurn int) {
	if yesno {
		term.MoveTo(tab, i+1)
		fmt.Print(sizeBar(maxSize, file.Size))
		fmt.Print(" ")
		tabTurn = tab + 13
	} else {
		tabTurn = tab
	}
	return tabTurn
}

package shko

import (
	"fmt"
	"strconv"
	"time"

	term "github.com/tj/go/term"
)

var (
	tab   = 20
	space = 1
)

func SelectInList(selected, scroll int, drawlist, childrens []File, parent File) {
	term.MoveTo(0, 0)
	term.ClearAll()
	termWidth, termHeight := term.Size()
	topSpace = 0
	sideSpace = 0
	if topBar {
		topSpace = 2
		Print(HighLight, Black, Cyan, DashBorder2(parent.Path, termWidth/2-(len([]rune(parent.Path)))/2))
		println()
	}
	if center && termHeight > len(drawlist) {
		topSpace += termHeight/2 - (len(drawlist) / 2)
		sideSpace = termWidth/2 - 25/2
	}
	if len(drawlist) == 0 {
		fmt.Print("  ")
		term.MoveTo(sideSpace+3, topSpace)
		Print(HighLight, Black, White, "  nothing to show  ")
	} else {
		var maxSize int64
		for _, el := range childrens {
			maxSize += el.Size
		}
		for i, el := range drawlist {
			if selected+scroll == el.Other.Number {
				el.Other.Active = true
			}
			fmt.Print("  ")
			if i == selected || el.Other.Selected == true {
				colorList(el, true, i+topSpace, maxSize)
			} else {
				colorList(el, false, i+topSpace, maxSize)
			}
			fmt.Print("\n")
			ResetStyle()
		}
	}
	if statBar {
		term.MoveTo(0, termHeight-2)
		println()
		Print(HighLight, Black, Cyan, DashBorder2(parent.Path, termWidth/2-(len([]rune(parent.Path)))/2))
	}
}

func colorList(file File, active bool, i int, maxSize int64) {
	termWidth, _ := term.Size()
	tab = space + 2 + sideSpace
	term.MoveTo(tab, i+1)
	if file.IsDir {
		if dirASwitch {
			Invert(active, HighLight, White)
		} else {
			Invert(active, HighLight, Green)
		}
	} else {
		Invert(active, Default, Cyan)
	}
	term.ClearLineEnd()
	tab = drawIcon(active, showIcons, file, i)
	tab = drawName(active, file, i)
	tab = drawChildren(showChildren, file, i)
	tab = drawMode(showMode, file, i)
	tab = drawDU(duMode, file, i, maxSize)
	tab = drawSize(showSize, file, i)
	tab = drawDate(showDate, file, i)
	//tab = drawMime(true, file, i)
	if !topBar {
		SetStyle(Default, White, Black)
	} else {
		term.MoveTo(termWidth-space, i+1)
		SetStyle(Default, White, Black)
	}
	term.ClearLineEnd()
}

func drawIcon(active, yesno bool, file File, i int) (tabTurn int) {
	before := ""
	after := "  "
	if file.Other.Selected && file.Other.Active {
		before = "×"
	} else if file.Other.Selected {
		before = " ×"
	} else if active {
		before = " "
	}
	before += " "
	if yesno {
		fmt.Print(before + file.Other.Icon + after)
	} else {
		if file.IsDir {
			fmt.Print(before + ">" + after)
		} else {
			fmt.Print(before + "-" + after)
		}
	}
	tabTurn = tab + 5
	return
}

func drawName(active bool, file File, i int) (tabTurn int) {
	term.MoveTo(tab, i+1)
	spacer := ""
	if active {
		spacer = " "
	}
	if file.IsDir {
		fmt.Print(spacer + file.Name + spacer + "/ ")
	} else {
		fmt.Print(spacer + file.Name + spacer + " ")
	}
	tabTurn = tab + 20
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
	if maxSize == 0 {
		maxSize = 1
	}
	percentage := int(size) * 100 / int(maxSize)
	for i := 1; i <= percentage; i = i + 10 {
		load += "█"
	}
	for i := 1; i <= 10-len([]rune(load)); i++ {
		uload += "░"
	}
	toPrint = "│" + load + uload + "│"
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

func drawMime(yesno bool, file File, i int) (tabTurn int) {
	if yesno {
		term.MoveTo(tab, i+1)
		fmt.Print(file.Mime)
		fmt.Print(" ")
		tabTurn = tab + 20
	} else {
		tabTurn = tab
	}
	return tabTurn
}

func statusWrite(toWrite string) {
	term.MoveTo(0, termHeight+1)
	Print(HighLight, Black, Red, toWrite)
}

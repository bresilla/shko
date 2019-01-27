package shko

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bresilla/dirk"
	term "github.com/bresilla/shko/term"
)

var (
	tab   = 20
	space = 1
)

func SelectInList(selected, scroll *int, drawlist, childrens *dirk.Files, currentDir *dirk.File) {
	term.MoveTo(0, 0)
	term.ClearAll()
	if showBar {
		Print(term.HighLight, term.Black, term.Cyan, DashBorder2(currentDir.Path, "-", termWidth/2-(len([]rune(currentDir.Path)))/2))
		Print(term.Default, term.Cyan, term.Black, DashBorder2("", "¯", 0))
	}
	if len(*drawlist) == 0 {
		term.MoveTo(termWidth/2-8, termHeight/2+1)
		Print(term.HighLight, term.Black, term.White, "  nothing to show  ")
	} else {
		for i := range *drawlist {
			if i == *selected || (*drawlist)[i].Selected == true {
				colorList((*drawlist)[i], true, i+barSpace)
			} else {
				colorList((*drawlist)[i], false, i+barSpace)
			}
			fmt.Print("\n")
			term.ResetStyle()
		}
	}
	if showBar {
		term.MoveTo(0, termHeight)
		Print(term.Default, term.Cyan, term.Black, DashBorder2("", "_", 0))
		Print(term.HighLight, term.Black, term.Cyan, DashBorder2(currentDir.Path, "-", termWidth/2-(len([]rune(currentDir.Path)))/2))
	}
}

func colorList(file *dirk.File, active bool, i int) {
	tab = space + 2 + sideSpace
	term.MoveTo(tab, i+1)
	if file.IsDir() {
		ColorSelect(active, term.HighLight, term.White)
	} else {
		ColorSelect(active, term.Default, term.Cyan)
	}
	tab = drawIcon(active, showIcons, file, i)
	tab = drawName(active, file, i)
	tab = drawChildren(showChildren, file, i)
	tab = drawMode(showMode, file, i)
	//tab = drawDU(dirk.DiskUse, file, i)
	tab = drawSize(showSize, file, i)
	tab = drawDate(showDate, file, i)
	tab = drawMime(showMime, file, i)
	term.SetStyle(term.Default, term.White, term.Black)
}

func ColorSelect(active bool, style term.Style, color term.Color) {
	if active {
		term.SetStyle(style, term.Black, color)
	} else {
		term.SetStyle(style, color, term.None)
	}
}

func drawIcon(active, yesno bool, file *dirk.File, i int) (tabTurn int) {
	before := ""
	after := "  "
	if file.Selected && file.Active {
		before = "×"
	} else if file.Selected {
		before = " ×"
	} else if active {
		before = " "
	}
	before += " "
	if yesno {
		fmt.Print(before + file.MimeIcon() + after)
	} else {
		if file.IsDir() {
			fmt.Print(before + ">" + after)
		} else {
			fmt.Print(before + "-" + after)
		}
	}
	tabTurn = tab + 5
	return
}

func drawName(active bool, file *dirk.File, i int) (tabTurn int) {
	term.MoveTo(tab, i+1)
	spacer := ""
	if active {
		spacer = " "
	}
	if file.IsDir() {
		fmt.Print(spacer + file.Nick + spacer + "/ ")
	} else {
		fmt.Print(spacer + file.Nick + spacer + " ")
	}
	tabTurn = tab + 20
	return
}

func drawChildren(yesno bool, file *dirk.File, i int) (tabTurn int) {
	if yesno {
		term.MoveTo(tab, i+1)
		fmt.Print("  " + strconv.Itoa(len(file.Childrens())) + " ")
		tabTurn = tab + 8
	} else {
		tabTurn = tab
	}
	return tabTurn
}

func drawSize(yesno bool, file *dirk.File, i int) (tabTurn int) {
	if yesno {
		term.MoveTo(tab, i+1)
		fmt.Print(file.SizeSTR(dirk.DiskUse) + " ")
		tabTurn = tab + 12
	} else {
		tabTurn = tab
	}
	return tabTurn
}

func drawDate(yesno bool, file *dirk.File, i int) (tabTurn int) {
	if yesno {
		term.MoveTo(tab, i+1)
		fmt.Print(file.TimeBirth().Format(time.RFC822) + " ")
		tabTurn = tab + 25
	} else {
		tabTurn = tab
	}
	return tabTurn
}

func drawMode(yesno bool, file *dirk.File, i int) (tabTurn int) {
	if yesno {
		term.MoveTo(tab, i+1)
		fmt.Print(file.File.Mode())
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

func drawDU(yesno bool, file *dirk.File, i int) (tabTurn int) {
	if yesno {
		term.MoveTo(tab, i+1)
		maxSize := file.MaxSize()
		fmt.Print(sizeBar(maxSize, file.SizeINT(dirk.DiskUse)))
		fmt.Print(" ")
		tabTurn = tab + 13
	} else {
		tabTurn = tab
	}
	return tabTurn
}

func drawMime(yesno bool, file *dirk.File, i int) (tabTurn int) {
	if yesno {
		term.MoveTo(tab, i+1)
		fmt.Print(file.MimeType())
		fmt.Print(" ")
		tabTurn = tab + 20
	} else {
		tabTurn = tab
	}
	return tabTurn
}

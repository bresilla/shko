package shko

import (
	"fmt"
	"log"
	"strconv"
	"time"

	. "./dirk"
	"github.com/logrusorgru/aurora"
	"github.com/peterh/liner"
	term "github.com/tj/go/term"
)

var (
	tab   = 20
	space = 1
	au    = aurora.NewAurora(colors)
)

func SelectInList(selected, scroll int, drawlist, childrens Files, currentDir File) {
	var maxSize int64
	lenMax := 15
	topSpace = 0
	sideSpace = 0

	for i := range childrens {
		if len(childrens[i].Name) > lenMax {
			lenMax = len(childrens[i].Name)
		}
		maxSize += childrens[i].Size
	}
	if center && termHeight > len(drawlist) {
		topSpace += termHeight/2 - (len(drawlist) / 2)
		sideSpace = termWidth/2 - lenMax/2 - 5
	}

	term.MoveTo(0, 0)
	term.ClearAll()
	if topBar {
		topSpace = 2
		Print(HighLight, Black, Cyan, DashBorder2(currentDir.Path, "-", termWidth/2-(len([]rune(currentDir.Path)))/2))
		Print(Default, Cyan, Black, DashBorder2("", "¯", 0))
	}
	if len(drawlist) == 0 {
		fmt.Print("  ")
		term.MoveTo(sideSpace+3, topSpace)
		Print(HighLight, Black, White, "  nothing to show  ")
	} else {
		for i, el := range drawlist {
			if selected+scroll == el.Number {
				el.Other.Active = true
			}
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
		term.MoveTo(0, termHeight)
		Print(Default, Cyan, Black, DashBorder2("", "_", 0))
		Print(HighLight, Black, Cyan, DashBorder2(currentDir.Path, "-", termWidth/2-(len([]rune(currentDir.Path)))/2))
	}
}

func colorList(file File, active bool, i int, maxSize int64) {
	tab = space + 2 + sideSpace
	term.MoveTo(tab, i+1)
	if file.IsDir {
		ColorSelect(active, HighLight, White)
	} else {
		ColorSelect(active, Default, Cyan)
	}
	tab = drawIcon(active, showIcons, file, i)
	tab = drawName(active, file, i)
	tab = drawChildren(showChildren, file, i)
	tab = drawMode(showMode, file, i)
	tab = drawDU(DiskUse, file, i, maxSize)
	tab = drawSize(showSize, file, i)
	tab = drawDate(showDate, file, i)
	tab = drawMime(showMime, file, i)
	SetStyle(Default, White, Black)
}

func ColorSelect(active bool, style Style, color Color) {
	if active {
		SetStyle(style, Black, color)
	} else {
		SetStyle(style, color, None)
	}
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
		fmt.Print("  " + strconv.Itoa(file.Other.Children) + " ")
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
		fmt.Print(file.BrtTime.Format(time.RFC822) + " ")
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
	cleanLine(0)
	term.MoveTo(0, termHeight+1)
	fmt.Print(au.Sprintf(au.BgCyan(au.Bold(au.Black(toWrite)))))
	fmt.Print(" ")
}

func statusRead(toWrite, defaultStr string) (text string) {
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

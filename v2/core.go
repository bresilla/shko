package shko

import (
	"fmt"
	"os"
	"strconv"

	"github.com/bresilla/dirk"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/rivo/tview"
)

func Main2() {

	dirk.IncHidden = true
	dirk.IncFiles = true
	childrenList := currentDir.ListDir()

	app := tview.NewApplication()
	pages := tview.NewPages()

	table := tview.NewTable().SetBorders(false)
	for i := range childrenList {
		color := tcell.ColorWhite
		if !childrenList[i].IsDir {
			color = tcell.ColorFuchsia
		}
		table.SetCell(i, 0, tview.NewTableCell(childrenList[i].Icon).SetTextColor(color).SetAlign(tview.AlignLeft))
		table.SetCell(i, 1, tview.NewTableCell(childrenList[i].Name).SetTextColor(color).SetAlign(tview.AlignLeft))
		table.SetCell(i, 2, tview.NewTableCell(strconv.Itoa(childrenList[i].Number)).SetTextColor(color).SetAlign(tview.AlignLeft))
	}
	table.Select(0, 0).SetSelectable(true, false).SetSelectedFunc(func(row int, column int) {
		table.GetCell(row, column).SetTextColor(tcell.ColorRed)
		table.SetSelectable(true, false)
	}).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
		if key == tcell.KeyEnter {
			table.SetSelectable(true, true)
			row, _ := table.GetSelection()
			if childrenList[row].IsDir {
				currentDir = childrenList[row]
			}
		}
	})

	pages.AddPage("main", table, true, true)
	if err := app.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}
}

func Main3() {
	dirk.IncHidden = true
	dirk.IncFiles = true
	childrenList := currentDir.ListDir()

	encoding.Register()

	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	screen.EnableMouse()

	pages := tview.NewPages()
	table := tview.NewTable().SetBorders(false)
	for i := range childrenList {
		color := tcell.ColorWhite
		if !childrenList[i].IsDir {
			color = tcell.ColorFuchsia
		}
		table.SetCell(i, 0, tview.NewTableCell(childrenList[i].Icon).SetTextColor(color).SetAlign(tview.AlignLeft))
		table.SetCell(i, 1, tview.NewTableCell(childrenList[i].Name).SetTextColor(color).SetAlign(tview.AlignLeft))
		table.SetCell(i, 2, tview.NewTableCell(strconv.Itoa(childrenList[i].Number)).SetTextColor(color).SetAlign(tview.AlignLeft))
	}
	pages.AddPage("main", table, true, true)
	pages.Draw(screen)

	for {
		screen.Show()
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				screen.Fini()
				os.Exit(0)
			} else if ev.Key() == tcell.KeyCtrlL {
				screen.Sync()
			} else {
				//print(ev.Rune())
			}
		case *tcell.EventMouse:
			button := ev.Buttons()
			button &= tcell.ButtonMask(0xff)
			if button != tcell.ButtonNone {
			}
			switch ev.Buttons() {
			case tcell.Button1:
			case tcell.Button2:
			case tcell.Button3:
			}
		}
	}
}

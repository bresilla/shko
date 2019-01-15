package shko

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/bresilla/dirk"
	"github.com/jroimartin/gocui"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/rivo/tview"
)

func handle(err error) {
	if err != nil {
		fmt.Print(err)
		log.Panicln(err)
	}
}

func Main() {

	dirk.IncHidden = true
	dirk.IncFiles = true
	currentDir, _ := dirk.MakeFile(os.Getenv("HOME"))
	childrenList := currentDir.ListDir()

	app := tview.NewApplication()
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
	table.Select(0, 0).SetDoneFunc(func(key tcell.Key) {
		table.SetSelectable(true, false)
	}).SetSelectedFunc(func(row int, column int) {
		table.GetCell(row, column).SetTextColor(tcell.ColorRed)
		table.SetSelectable(true, false)
	})
	if err := app.SetRoot(table, true).Run(); err != nil {
		panic(err)
	}
}

func Main1() {
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

	screen.SetStyle(tcell.StyleDefault.Background(tcell.ColorRed).Foreground(tcell.ColorWhite))
	screen.EnableMouse()
	screen.Show()
	screen.Clear()

	for {
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
				print(ev.Rune())
			}
		case *tcell.EventMouse:
			button := ev.Buttons()
			if button&tcell.WheelUp != 0 {
				print("WheelUp")
			}
			if button&tcell.WheelDown != 0 {
				print("WheelDown")
			}
			if button&tcell.WheelLeft != 0 {
				print("WheelLeft")
			}
			if button&tcell.WheelRight != 0 {
				print("WheelRight")
			}
			button &= tcell.ButtonMask(0xff)
			if button != tcell.ButtonNone {
				x, y := ev.Position()
				print(x, y)
			}
			switch ev.Buttons() {
			case tcell.Button1:
				print('1')
			case tcell.Button2:
				print('2')
			case tcell.Button3:
				print('3')
			}
		}
	}
}

func Main4() {
	gui, err := gocui.NewGui(gocui.OutputNormal)
	handle(err)

	gui.SetManagerFunc(layout)

	if err := gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(gui *gocui.Gui) error {
	maxX, maxY := gui.Size()
	if v, err := gui.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Hello world!")
	}
	return nil
}

func quit(gui *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

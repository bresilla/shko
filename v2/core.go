package shko

import (
	"fmt"
	"log"
	"os"
	"regexp"

	dirk "github.com/bresilla/dirk"
	"github.com/jroimartin/gocui"
	component "github.com/skanehira/gocui-component"
)

var (
	inp         *component.InputField
	startDir, _ = dirk.MakeFile(os.Getenv("PWD"))
)

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
func otherview(g *gocui.Gui, v *gocui.View) error {
	g.SetViewOnTop("v1")
	return nil
}
func overview(g *gocui.Gui, v *gocui.View) error {
	g.SetViewOnTop("v2")
	return nil
}
func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'p', gocui.ModNone, otherview); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'o', gocui.ModNone, overview); err != nil {
		return err
	}

	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("v1", -1, -1, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, "Item 1")
		fmt.Fprintln(v, "Item 2")
	}
	if v, err := g.SetView("v2", -1, -1, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		inp.Draw()
	}
	return nil
}

func startNumber(value string) bool {
	return regexp.MustCompile(`^[0-9]`).MatchString(value)
}

func endString(value string) bool {
	return regexp.MustCompile(`[a-zA-Z]$`).MatchString(value)
}

func Main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.Cursor = true
	g.SetManagerFunc(layout)
	inp = component.NewInputField(g, "password", 0, 0, 10, 15).SetText("trim")
	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

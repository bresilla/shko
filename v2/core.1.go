package shko

import (
	"fmt"
	"log"
	"os"

	"github.com/bresilla/dirk"
	"github.com/jroimartin/gocui"
)

func handle(err error) {
	if err != nil {
		fmt.Print(err)
		log.Panicln(err)
	}
}

var (
	currentDir, _ = dirk.MakeFile(os.Getenv("HOME"))
)

func Main() {

	dirk.IncHidden = true
	dirk.IncFiles = true
	//childrenList := currentDir.ListDir()

	gui, err := gocui.NewGui(gocui.OutputNormal)
	handle(err)
	layout := 1
	if err := manager(layout, gui); err != nil {
		log.Print(err)
	}

	if err := gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func manager(lay int, g *gocui.Gui) error {
	switch lay {
	case 1:
		g.SetManagerFunc(layout_1)
	case 2:
		g.SetManagerFunc(layout_1)
	}
	return nil
}

func layout_1(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if _, err := g.SetView("main", -1, -1, maxX, maxY-2); err != nil &&
		err != gocui.ErrUnknownView {
		return err
	}
	if _, err := g.SetView("status", -1, maxY-2, maxX, maxY); err != nil &&
		err != gocui.ErrUnknownView {
		return err
	}
	return nil
}

func keybindings(g *gocui.Gui) error {
	err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	})
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", '1', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		_, err := g.SetViewOnTop("v1")
		return err
	})
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", '2', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		_, err := g.SetViewOnTop("v2")
		return err
	})
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", '3', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		_, err := g.SetViewOnTop("v3")
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

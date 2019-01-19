package shko

import (
	"fmt"
	"log"
	"os"

	"github.com/bresilla/dirk"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	runewidth "github.com/mattn/go-runewidth"
)

func handle(err error) {
	if err != nil {
		fmt.Print(err)
		log.Panicln(err)
	}
}

var (
	currentDir, _ = dirk.MakeFile(os.Getenv("HOME"))
	childrenList  = currentDir.ListDir()
	layout        = 1
)

var defStyle tcell.Style

func WriteStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}

func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, r rune) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, style)
		s.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, style)
		s.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}
	if y1 != y2 && x1 != x2 {
		// Only add corners if we need to
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		for col := x1 + 1; col < x2; col++ {
			s.SetContent(col, row, r, nil, style)
		}
	}
}

func Main() {

	encoding.Register()

	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	defStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	s.SetStyle(defStyle)
	s.EnableMouse()
	s.Clear()

	white := tcell.StyleDefault.
		Foreground(tcell.ColorWhite).Background(tcell.ColorRed)

	w, h := s.Size()

	for {
		drawBox(s, 1, 1, 42, 6, white, ' ')
		drawBox(s, w-42, h-6, w-1, h-1, white, ' ')
		WriteStr(s, 2, 2, white, "Press ESC twice to exit, C to clear.")
		s.Show()

		ev := s.PollEvent()
		w, h = s.Size()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Clear()
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Rune() == 'q' {
				s.Fini()
				os.Exit(0)
			} else if ev.Key() == tcell.KeyCtrlL {
				s.Sync()
			}
		case *tcell.EventMouse:
			button := ev.Buttons()
			button &= tcell.ButtonMask(0xff)
			switch ev.Buttons() {
			case tcell.Button1:
			case tcell.Button2:
			case tcell.Button3:
			case tcell.Button4:
			case tcell.Button5:
			case tcell.Button6:
			case tcell.Button7:
			case tcell.Button8:
			}

		}
	}
}

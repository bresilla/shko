package spejt

import (
	"fmt"

	term "github.com/nsf/termbox-go"
)

func Run() {
	err := term.Init()
	ErrorCheck(err)
	defer term.Close()

keyPressListenerLoop:
	for {
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			if ev.Key == term.KeySpace {
				fmt.Println("SPACEEEEE")
			} else if ev.Key == term.KeyEsc || ev.Ch == 45 || ev.Ch == 113 {
				break keyPressListenerLoop
			} else if ev.Key == term.KeyArrowUp {
				term.Sync()
				fmt.Println(term.KeyArrowUp)
			} else if ev.Key == term.KeyArrowDown {
				term.Sync()
				fmt.Println(term.KeyArrowDown)
			} else if ev.Key == term.KeyArrowLeft {
				term.Sync()
				fmt.Println(term.KeyArrowLeft)
			} else if ev.Key == term.KeyArrowRight {
				term.Sync()
				fmt.Println(term.KeyArrowRight)
			} else {
				term.Sync()
				fmt.Println("ASCII : ", ev.Ch)
			}
		case term.EventError:
			panic(ev.Err)
		}
	}
}

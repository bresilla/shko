package spejt

import (
	"fmt"

	term "github.com/nsf/termbox-go"
)

func reset() {
	term.Sync()
}

func Run() {
	err := term.Init()
	ErrorCheck(err)
	defer term.Close()

keyPressListenerLoop:
	for {
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			switch ev.Key {
			case term.KeyEsc | term.KeySpace:
				break keyPressListenerLoop
			case term.KeyArrowUp:
				reset()
				fmt.Println("Arrow Up pressed")
			case term.KeyArrowDown:
				reset()
				fmt.Println("Arrow Down pressed")
			case term.KeyArrowLeft:
				reset()
				fmt.Println("Arrow Left pressed")
			case term.KeyArrowRight:
				reset()
				fmt.Println("Arrow Right pressed")
			case term.KeyEnter:
				reset()
				fmt.Println("Enter pressed")
			default:
				break keyPressListenerLoop
			}
		case term.EventError:
			panic(ev.Err)
		}
	}
}

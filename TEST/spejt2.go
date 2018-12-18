package spejt

import (
	"os"

	term "github.com/nsf/termbox-go"
)

func Run() {
	err := term.Init()
	ErrorCheck(err)
	defer term.Close()

keyPressListenerLoop:
	for {
		children, parent := ListDirs(currentDir)
		SelectInList(number, children)
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			if ev.Key == term.KeyEsc || ev.Ch == 45 {
				break keyPressListenerLoop
			} else if ev.Key == term.KeyEnter || ev.Key == term.KeyCtrlC {
				break keyPressListenerLoop
			} else if ev.Key == term.KeyArrowUp {
				term.Sync()
				number--
				if number < 0 {
					if cursoroll {
						number = len(children) - 1
					} else {
						number = 0
					}
				}
			} else if ev.Key == term.KeyArrowDown {
				term.Sync()
				number++
				if number > len(children)-1 {
					if cursoroll {
						number = 0
					} else {
						number = len(children) - 1
					}
				}
			} else if ev.Key == term.KeyArrowLeft {
				term.Sync()
				if len(cursorarr) > 0 {
					number = cursorarr[len(cursorarr)-1]
					cursorarr = cursorarr[:len(cursorarr)-1]
				} else {
					number = 0
				}
				currentDir = CurrentDir(parent.Path)
			} else if ev.Key == term.KeyArrowRight {
				term.Sync()
				if children[number].IsDir {
					cursorarr = append(cursorarr, number)
					currentDir = CurrentDir(children[number].Path)
				} else {
					OpenFile(children[number])
				}
				number = 0
			} else {
				term.Sync()
				//fmt.Println("ASCII : ", ev.Ch)
			}
		case term.EventError:
			panic(ev.Err)
		}
	}
	file, _ := os.Create(outDir)
	file.WriteString(currentDir.Path)
	term.Flush()
}

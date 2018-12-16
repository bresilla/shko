package spejt

import (
	"fmt"
	"os"

	term "github.com/nsf/termbox-go"
)

var (
	incFolder = true
	incFiles  = true
	incHidden = false
)

var (
	currentDir = StringDirToFile(os.Getenv("PWD"))
	parent     string
	child      string
)

func MakeCool(file []File) {
	for _, f := range file {
		fmt.Println(f.Name)
	}
}

func Run() {
	err := term.Init()
	ErrorCheck(err)
	defer term.Close()
	children, _ := ListDirs(currentDir)
	MakeCool(children)

keyPressListenerLoop:
	for {
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			if ev.Key == term.KeySpace {
				currentDir = CurrentDir("/home/bresilla/DATA")
				fmt.Println(currentDir.Path, currentDir.Parent)
				//StatDir("/")
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
				_, parent := ListDirs(currentDir)
				children, _ := ListDirs(CurrentDir(parent.Path))
				MakeCool(children)
				currentDir = CurrentDir(parent.Path)
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

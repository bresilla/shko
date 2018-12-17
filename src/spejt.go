package spejt

import (
	"fmt"
	"os"
	"os/exec"

	term "github.com/nsf/termbox-go"
)

var (
	incFolder       = true
	incFiles        = true
	incHidden       = false
	cursoroll       = true
	cursorset       = false
	shortcut   rune = 113
	currentDir      = StringDirToFile(os.Getenv("PWD"))
	parent     string
	child      string
	cmd        *exec.Cmd
	startNr    = 0
	cursorarr  []int
)

func ChangeDir(cmd *exec.Cmd, filepath string) {
	cmd = exec.Command("cd", filepath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
}

func Run() {
	err := term.Init()
	ErrorCheck(err)
	defer term.Close()

keyPressListenerLoop:
	for {
		_, parent := ListDirs(currentDir)
		children, _ := ListDirs(currentDir)
		SelectInList(startNr, children)
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			if ev.Ch == shortcut {
				break keyPressListenerLoop
			} else if ev.Key == term.KeyEsc || ev.Ch == 45 {
				break keyPressListenerLoop
			} else if ev.Key == term.KeyArrowUp {
				term.Sync()
				startNr--
				if startNr < 0 {
					if cursoroll {
						startNr = len(children) - 1
					} else {
						startNr = 0
					}
				}
			} else if ev.Key == term.KeyArrowDown {
				term.Sync()
				startNr++
				if startNr > len(children)-1 {
					if cursoroll {
						startNr = 0
					} else {
						startNr = len(children) - 1
					}
				}
			} else if ev.Key == term.KeyArrowLeft {
				term.Sync()
				if len(cursorarr) > 0 && !cursorset {
					startNr = cursorarr[len(cursorarr)-1]
					cursorarr = cursorarr[:len(cursorarr)-1]
				} else {
					startNr = 0
				}
				currentDir = CurrentDir(parent.Path)
			} else if ev.Key == term.KeyArrowRight {
				cursorarr = append(cursorarr, startNr)
				term.Sync()
				if children[startNr].IsDir {
					currentDir = CurrentDir(children[startNr].Path)
				} else {
					OpenInEditor(children[startNr].Path)
				}
				startNr = 0
			} else {
				term.Sync()
				fmt.Println("ASCII : ", ev.Ch)
			}
		case term.EventError:
			panic(ev.Err)
		}
	}
}

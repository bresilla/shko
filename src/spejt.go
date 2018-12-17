package spejt

import (
	"fmt"
	"os"
	"os/exec"

	term "github.com/nsf/termbox-go"
)

var (
	incFolder  = true
	incFiles   = true
	incHidden  = true
	currentDir = StringDirToFile(os.Getenv("PWD"))
	parent     string
	child      string
	cmd        *exec.Cmd
	startNr    = 0
)

func MakeCool(file []File) {
	for _, f := range file {
		fmt.Println(f.Name)
	}
}

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
			if ev.Key == term.KeySpace {
				fmt.Println("SPACE")
			} else if ev.Key == term.KeyEsc || ev.Ch == 45 {
				break keyPressListenerLoop
			} else if ev.Key == term.KeyArrowUp {
				term.Sync()
				startNr--
				if startNr < 0 {
					startNr = len(children) - 1
				}
			} else if ev.Key == term.KeyArrowDown {
				term.Sync()
				startNr++
				if startNr > len(children)-1 {
					startNr = 0
				}
			} else if ev.Key == term.KeyArrowLeft {
				term.Sync()
				startNr = 0
				currentDir = CurrentDir(parent.Path)
			} else if ev.Key == term.KeyArrowRight {
				term.Sync()
			} else {
				term.Sync()
				//fmt.Println("ASCII : ", ev.Ch)
			}
		case term.EventError:
			panic(ev.Err)
		}
	}
}

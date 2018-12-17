package spejt

import (
	"fmt"
	"os"
	"os/exec"

	term "github.com/buger/goterm"
)

var (
	incFolder  = true
	incFiles   = true
	incHidden  = false
	cursoroll  = true
	cursorset  = false
	shortcut   = 113
	currentDir = StringDirToFile(os.Getenv("PWD"))
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

func Loop() {
	term.MoveCursor(0, 0)
	term.Clear()
	term.Flush()
	for {
		term.MoveCursor(0, 0)
		_, parent := ListDirs(currentDir)
		children, _ := ListDirs(currentDir)
		SelectInList(startNr, children)
		ascii, keycode, _ := GetChar()
		if ascii == 3 || ascii == 45 {
			break
		} else if ascii == shortcut {
			break
		} else if keycode == 38 {
			//up
			startNr--
			if startNr < 0 {
				if cursoroll {
					startNr = len(children) - 1
				} else {
					startNr = 0
				}
			}
		} else if keycode == 40 {
			//down
			startNr++
			if startNr > len(children)-1 {
				if cursoroll {
					startNr = 0
				} else {
					startNr = len(children) - 1
				}
			}
		} else if keycode == 37 {
			//left
			if len(cursorarr) > 0 && !cursorset {
				startNr = cursorarr[len(cursorarr)-1]
				cursorarr = cursorarr[:len(cursorarr)-1]
			} else {
				startNr = 0
			}
			currentDir = CurrentDir(parent.Path)
		} else if keycode == 39 {
			//right
			cursorarr = append(cursorarr, startNr)
			if children[startNr].IsDir {
				currentDir = CurrentDir(children[startNr].Path)
			} else {
				OpenInEditor(children[startNr].Path)
			}
			startNr = 0
		} else {
			fmt.Println(ascii, "\t", keycode)
		}
		term.Flush()
		term.Clear()
	}
	fmt.Println()
}

package spejt

import (
	"bufio"
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
	outDir     = "/tmp/spejt"
)

func OpenItem(startNr int, children []File, currentDir File, cursorarr *[]int) (file File) {
	if children[startNr].IsDir {
		*cursorarr = append(*cursorarr, startNr)
		file = CurrentDir(children[startNr].Path)
	} else {
		OpenInEditor(children[startNr].Path)
		file = currentDir
	}
	return
}

func WriteTempFile(toWrite string) {
	file, err := os.OpenFile(outDir, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("File does not exists or cannot be created")
		os.Exit(1)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	fmt.Fprintf(w, toWrite)
	w.Flush()
}

func AlterWrite(toWrite string) {
}

func Loop() {
	term.MoveCursor(0, 0)
	term.Flush()
	term.Clear()
	for {
		children, parent := ListDirs(currentDir)
		term.MoveCursor(0, 0)
		term.Flush()
		term.Clear()
		fmt.Println()
		SelectInList(startNr, children)
		fmt.Println()
		ascii, keycode, _ := GetChar()
		if ascii == 3 || ascii == 45 || ascii == 13 {
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
			currentDir = OpenItem(startNr, children, currentDir, &cursorarr)
			startNr = 0
		} else {
			fmt.Println(ascii, "\t", keycode)
		}
		fmt.Println()
	}
	fmt.Println()
	file, _ := os.Create(outDir)
	file.WriteString(currentDir.Path)
}

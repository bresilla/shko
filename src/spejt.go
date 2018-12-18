package spejt

import (
	"fmt"
	"os"

	term "github.com/buger/goterm"
)

var (
	incFolder     = true
	incFiles      = true
	incHidden     = false
	cursoroll     = true
	shortcut      = 113
	currentDir, _ = makeFile(os.Getenv("PWD"))
	parent        string
	child         string
	number        = 0
	cursorarr     []int
	outDir        = "/tmp/spejt"
)

func Loop() {
	term.MoveCursor(0, 0)
	fmt.Print("\033[?25l")
	term.Flush()
	term.Clear()
	for {
		subdirs, parent := ListDirs(currentDir)
		term.MoveCursor(0, 0)
		term.Flush()
		term.Clear()
		fmt.Println()
		SelectInList(number, subdirs)
		ascii, keycode, _ := GetChar()
		if ascii == 3 || ascii == 27 || ascii == 13 {
			break
		} else if ascii == shortcut {
			break
		} else if keycode == 38 {
			//up
			number--
			if number < 0 {
				if cursoroll {
					number = len(subdirs) - 1
				} else {
					number = 0
				}
			}
		} else if keycode == 40 {
			//down
			number++
			if number > len(subdirs)-1 {
				if cursoroll {
					number = 0
				} else {
					number = len(subdirs) - 1
				}
			}
		} else if keycode == 37 {
			//left
			if len(cursorarr) > 0 {
				number = cursorarr[len(cursorarr)-1]
				cursorarr = cursorarr[:len(cursorarr)-1]
			} else {
				number = 0
			}
			currentDir, _ = makeFile(parent.Path)
		} else if keycode == 39 {
			//right
			if len(subdirs) == 0 {
				continue
			}
			if subdirs[number].IsDir {
				cursorarr = append(cursorarr, number)
				currentDir, _ = makeFile(subdirs[number].Path)
			} else {
				OpenFile(subdirs[number])
			}
			number = 0
		} else {
			for {
				fmt.Println()
				fmt.Println(len(subdirs))
				//file, _ := ListRecourPathsNFiles(outDir)
				//println(file)
				ascii, keycode, _ := GetChar()
				break
				if ascii == 3 || keycode == 50 {
					break
				} else {
					//fmt.Println(ascii, "\t", keycode)
				}
			}
		}
		fmt.Println()
		fmt.Print("\033[?25l")
	}
	fmt.Print("\033[?25h")
	fmt.Println()
	file, _ := os.Create(outDir)
	file.WriteString(currentDir.Path)
}

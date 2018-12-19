package spejt

import (
	"fmt"
	"os"

	term "github.com/buger/goterm"
)

var (
	incFolder     = true
	incFiles      = false
	incHidden     = false
	cursoroll     = true
	shortcut      = 113
	currentDir, _ = makeFile(os.Getenv("PWD"))
	foreward      = true
	backward      = true
	number        = 0
	cursorarr     []int
	outDir        = "/tmp/spejt"
	track         = 0
)

func Loop() {
	fmt.Print("\033[?25l")
	term.Flush()
	term.Clear()
	for {
		children, parent := ListDirs(currentDir)
		subdirs := children
		if len(subdirs) > term.Height()-1 {
			if track < 0 {
				track = 0
			} else if track > len(children)+1-term.Height() {
				track = len(children) + 1 - term.Height()
			}
			subdirs = subdirs[0+track : term.Height()-1+track]
		}
		SelectInList(number, subdirs)
		ascii, keycode, _ := GetChar()
		//println(ascii, "\t", keycode)
		if ascii == 3 || ascii == 27 || ascii == 13 {
			break
		} else if ascii == shortcut || keycode == shortcut {
			break
		} else if keycode == 38 { //up
			if backward {
				number--
			} else {
				track--
			}
			if number < 0 {
				if cursoroll {
					number = len(subdirs) - 1
				} else {
					number = 0
				}
			}
		} else if keycode == 40 { //down
			if foreward {
				number++
			} else {
				track++
			}
			if number > len(subdirs)-1 {
				if cursoroll {
					number = 0
				} else {
					number = len(subdirs) - 1
				}
			}
		} else if keycode == 37 { //left
			if len(cursorarr) > 0 {
				number = cursorarr[len(cursorarr)-1]
				cursorarr = cursorarr[:len(cursorarr)-1]
			} else {
				number = 0
			}
			currentDir, _ = makeFile(parent.Path)
		} else if keycode == 39 { //right
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
				ascii, keycode, _ := GetChar()
				if ascii == 46 {
					if incHidden {
						incHidden = false
					} else {
						incHidden = true
					}
					break
				} else if ascii == 44 {
					if incFiles {
						incFiles = false
					} else {
						incFiles = true
					}
					break
				} else if ascii == 35 {
					if cursoroll {
						cursoroll = false
					} else {
						cursoroll = true
					}
					break
				} else if ascii == 3 || keycode == 50 {
					break
				} else if ascii == 45 {
					track--
					break
				} else {
					track++
					break
				}
			}
		}
	}
	fmt.Print("\033[?25h")
	fmt.Println()
	file, _ := os.Create(outDir)
	file.WriteString(currentDir.Path)
}

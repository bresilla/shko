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
	wrap          = true
	shortcut      = 113
	currentDir, _ = makeFile(os.Getenv("PWD"))
	startDir, _   = makeFile(os.Getenv("PWD"))
	changeDir     = true
	number        = 0
	scroll        = 0
	cursorarr     []int
	outDir        = "/tmp/spejt"
	file, _       = os.Create(outDir)
)

func Loop() {
	fmt.Print("\033[?25l")
	term.Flush()
	term.Clear()
	children, parent := ListDirs(currentDir)
	for {
		var foreward = false
		var backward = false
		subdirs := children
		if len(subdirs) > term.Height()-1 {
			if number > term.Height()/2 {
				foreward = true
				backward = false
			} else if number < term.Height()/2-2 {
				backward = true
				foreward = false
			}
			if scroll <= 0 {
				scroll = 0
				backward = false
			} else if scroll >= len(children)+1-term.Height() {
				scroll = len(children) + 1 - term.Height()
				foreward = false
			}
			subdirs = subdirs[0+scroll : term.Height()-1+scroll]
		}
		SelectInList(number, subdirs)
		ascii, keycode, _ := GetChar()
		if ascii == 13 || ascii == shortcut || keycode == shortcut {
			break
		} else if ascii == 27 || ascii == 3 {
			changeDir = false
			break
		} else if keycode == 38 { //up
			if backward {
				scroll--
			} else {
				number--
			}
			if number < 0 {
				if wrap {
					number = len(subdirs) - 1
					scroll = len(children) - 1
				} else {
					number = 0
				}
			}
		} else if keycode == 40 { //down
			if foreward {
				scroll++
			} else {
				number++
			}
			if number > len(subdirs)-1 {
				if wrap {
					number = 0
					scroll = 0
				} else {
					number = len(subdirs) - 1
				}
			}
		} else if keycode == 37 { //left
			backward = false
			foreward = false
			oldDir := currentDir
			currentDir, _ = makeFile(parent.Path)
			children, parent = ListDirs(currentDir)
			number, scroll = find(children, oldDir)
		} else if keycode == 39 { //right
			if len(subdirs) == 0 {
				continue
			}
			if subdirs[number].IsDir {
				currentDir, _ = makeFile(subdirs[number].Path)
				children, parent = ListDirs(currentDir)
			} else {
				OpenFile(subdirs[number])
				fmt.Print("\033[?25l")
			}
			number = 0
			scroll = 0
			backward = false
			foreward = false
		} else {
			for {
				ascii, keycode, _ := GetChar()
				if ascii == 46 {
					incHidden = !incHidden
					children, parent = ListDirs(currentDir)
					break
				} else if ascii == 44 {
					incFiles = !incFiles
					children, parent = ListDirs(currentDir)
					break
				} else if ascii == 35 {
					wrap = !wrap
					children, parent = ListDirs(currentDir)
					break
				} else {
					println(ascii, "\t", keycode)
					GetChar()
					break
				}
			}
		}
	}
	fmt.Print("\033[?25h")
	fmt.Println()
	if changeDir {
		file.WriteString(currentDir.Path)
	} else {
		file.WriteString(startDir.Path)
	}
}
